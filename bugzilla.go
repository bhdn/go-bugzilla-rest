// Package bugzilla can get bugs, attachments and update them using Bugzilla's
// REST API. This is based on an old code that used the web interface to update
// bugs, so the interface may not be as nice as it could be.
package bugzilla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// RequestError happens when building the request
type RequestError struct{ error }

func (e RequestError) Error() string {
	return fmt.Sprintf("cannot build request: %v", e.error)
}

// BugzillaError results from structured Bugzilla errors
type BugzillaError struct{ error }

func (e BugzillaError) Error() string {
	return fmt.Sprintf("Bugzilla: %v", e.error)
}

// ConnectionError happens when performing the request
type ConnectionError struct{ error }

func (e ConnectionError) Error() string {
	return fmt.Sprintf("cannot communicate with server: %v", e.error)
}

// DecodeErrror happens when it fails to parse JSON from the responses
type DecodeErrror struct{ error }

func (e DecodeErrror) Error() string {
	return fmt.Sprintf("Error decoding response: %v", e.error)
}

// Cacher should be anything that takes the name of the object to be cached
// and returns something that can receive writes with the contents and then
// eventually be closed.
type Cacher interface {
	GetWriter(id string) io.WriteCloser
}

// Config sets the parameters needed to set up the client. Cacher can be
// left zeroed.
type Config struct {
	BaseURL  string
	Username string
	ApiKey   string
	Cacher   Cacher
}

func (c *Config) emailAddress() (string, error) {
	// Consider adding .Email later
	email := c.Username
	if !strings.Contains(email, "@") {
		return "", ErrBugzilla{fmt.Errorf("Your username doesn't look like an email address: %v", email)}
	}
	return email, nil
}

// Client keeps the state of the client.
type Client struct {
	Config        Config
	seriousClient *http.Client
	cacher        Cacher
}

func getHTTPClient(config *Config) *http.Client {
	tr := http.DefaultClient.Transport
	client := http.Client{Transport: tr}
	rt := useHeader(tr)
	rt.Set("Content-type", "application/json")
	rt.Set("Accept", "application/json")
	client.Transport = rt
	return &client
}

// New prepares a *Client for connecting to the Bugzilla Web interface
func New(config Config) (*Client, error) {
	httpClient := getHTTPClient(&config)
	client := &Client{Config: config, seriousClient: httpClient, cacher: config.Cacher}
	return client, nil
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func useHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

func (c *Client) addAuthOptions(values *url.Values) {
	values.Add("Bugzilla_api_key", c.Config.ApiKey)
}

func (c *Client) makeURL(base_path string, values *url.Values) (string, error) {
	url, err := url.Parse(c.Config.BaseURL)
	if err != nil {
		return "", RequestError{err}
	}
	c.addAuthOptions(values)
	url.RawQuery = values.Encode()
	url.Path = path.Join(url.Path, base_path)

	return url.String(), nil
}

func valuesFromBugIds(ids []int) (first int, values *url.Values, err error) {
	values = &url.Values{}
	if len(ids) < 1 {
		err = fmt.Errorf("Not enough IDs passed to getCommentsURL")
		return
	}
	for _, id := range ids {
		values.Add("ids", fmt.Sprintf("%d", id))
	}
	first = ids[0]
	return
}

func (c *Client) valuesFromMap(m map[string]string) *url.Values {
	values := &url.Values{}
	for k, v := range m {
		values.Add(k, v)
	}
	return values
}

func (c *Client) getCommentsURL(bugIds []int, params map[string]string) (string, error) {
	first, values, err := valuesFromBugIds(bugIds)
	if err != nil {
		return "", err
	}
	for k, v := range params {
		values.Add(k, v)
	}
	path := fmt.Sprintf("/rest/bug/%d/comment", first)
	return c.makeURL(path, values)
}

func (c *Client) getAttachmentsURL(bugIds []int, params map[string]string) (string, error) {
	first, values, err := valuesFromBugIds(bugIds)
	if err != nil {
		return "", err
	}
	for k, v := range params {
		values.Add(k, v)
	}
	path := fmt.Sprintf("/rest/bug/%d/attachment", first)
	return c.makeURL(path, values)
}

func (c *Client) getAttachmentURL(attachmentId int, params map[string]string) (string, error) {
	path := fmt.Sprintf("/rest/bug/attachment/%d", attachmentId)
	return c.makeURL(path, c.valuesFromMap(params))
}

func (c *Client) getBugsURL(ids []int, params map[string]string) (string, error) {
	//values := &url.Values{}
	first, values, err := valuesFromBugIds(ids)
	if err != nil {
		return "", err
	}
	for k, v := range params {
		values.Add(k, v)
	}
	path := fmt.Sprintf("/rest/bug/%d", first)
	return c.makeURL(path, values)
}

func (c *Client) getUpdateBugsURL(ids []int, params map[string]string) (string, error) {
	return c.getBugsURL(ids, params)
}

func (c *Client) getDownloadURL(id int) (string, error) {
	url, err := url.Parse(c.Config.BaseURL)
	if err != nil {
		return "", RequestError{err}
	}

	url.Path = path.Join(url.Path, "attachment.cgi")

	query := url.Query()
	query.Set("id", fmt.Sprintf("%d", id))
	url.RawQuery = query.Encode()

	return url.String(), nil
}

// User represents user as used in assigned_to, comment author and other
// fields (except Cc.)
type User struct {
	Id       int    `json:"id"`
	Name     string `xml:"name,attr" json:"name"`
	Email    string `xml:",chardata" json:"email"`
	RealName string `json:"real_name"`
}

// Flag represents flags such as needinfo
type Flag struct {
	Name             string    `xml:"name,attr" json:"name"`
	ID               int       `xml:"id,attr" json:"id"`
	TypeID           int       `xml:"type_id,attr" json:"type_id"`
	Status           string    `xml:"status,attr" json:"status"`
	Setter           string    `xml:"setter,attr" json:"setter"`
	Requestee        string    `xml:"requestee,attr" json:"requestee"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
}

type Bug struct {
	ActualTime          float64    `json:"actual_time"`
	Alias               []string   `json:"alias"`
	AssignedTo          string     `json:"assigned_to"`
	AssignedToDetail    User       `json:"assigned_to_detail"`
	Blocks              []int      `json:"blocks"`
	CC                  []string   `json:"cc"`
	CCDetail            []User     `json:"cc_detail"`
	Classification      string     `json:"classification"`
	Component           string     `json:"component"`
	CreationTime        time.Time  `json:"creation_time"`
	Creator             string     `json:"creator"`
	CreatorDetail       User       `json:"creator_detail"`
	Deadline            *time.Time `json:"deadline"`
	DependsOn           []int      `json:"depends_on"`
	DupeOf              *int       `json:"dupe_of"`
	EstimatedTime       float64    `json:"estimated_time"`
	Flags               []Flag     `json:"flags"`
	Groups              []string   `json:"groups"`
	ID                  int        `json:"id"`
	IsCCAccessible      bool       `json:"is_cc_accessible"`
	IsConfirmed         bool       `json:"is_confirmed"`
	IsCreatorAccessible bool       `json:"is_creator_accessible"`
	IsOpen              bool       `json:"is_open"`
	Keywords            []string   `json:"keywords"`
	LastChangeTime      time.Time  `json:"last_change_time"`
	OpSys               string     `json:"op_sys"`
	Platform            string     `json:"platform"`
	Priority            string     `json:"priority"`
	Product             string     `json:"product"`
	QAContact           string     `json:"qa_contact"`
	QAContactDetail     User       `json:"qa_contact_detail"`
	RemainingTime       float64    `json:"remaining_time"`
	Resolution          string     `json:"resolution"`
	SeeAlso             []string   `json:"see_also"`
	Severity            string     `json:"severity"`
	Status              string     `json:"status"`
	Summary             string     `json:"summary"`
	TargetMilestone     string     `json:"target_milestone"`
	UpdateToken         string     `json:"update_token"`
	URL                 string     `json:"url"`
	Version             string     `json:"version"`
	Whiteboard          string     `json:"whiteboard"`

	Attachments []Attachment `xml:"attachment" json:"attachments"`
	Comments    []Comment    `xml:"long_desc" json:"comments"`
}

// hasNeedinfoFor check if a given email has been set to needinfo already
func (b *Bug) hasNeedinfoFor(email string) bool {
	return len(b.findNeedinfosFor(email)) > 0
}

// findNeedinfosFor gets the ids of the needinfos for a given email
func (b *Bug) findNeedinfosFor(email string) []int {
	ids := make([]int, 0, 0)
	lowerEmail := strings.ToLower(email)
	for _, flag := range b.Flags {
		if flag.Name == "needinfo" {
			if lowerEmail == "" || strings.ToLower(flag.Requestee) == lowerEmail {
				ids = append(ids, flag.ID)
			}
		}
	}
	return ids
}

// Attachment as provided by the bug information page. This struct has only
// name, size and attachid set when coming from DownloadAttachment, as it's
// extracted from the HTTP headers.
type Attachment struct {
	ID    int `json:"id,omitempty"`
	BugId int `json:"bug_id,omitempty"`

	Creator        string    `json:"creator,omitempty"`
	Data           []byte    `json:"data,omitempty"`
	IsObsolete     int       `json:"is_obsolete,omitempty"`
	IsPatch        int       `json:"is_patch,omitempty"`
	IsPrivate      int       `json:"is_private,omitempty"`
	CreationTime   time.Time `json:"creation_time,omitempty"`
	LastChangeTime time.Time `json:"last_change_time,omitempty"`
	ContentType    string    `json:"content_type,omitempty"`
	Summary        string    `json:"summary,omitempty"`
	Filename       string    `json:"file_name,omitempty"`
	Size           int64     `json:"size,omitempty"`
	Token          string    `json:"token,omitempty"`

	Flags []Flag `json:"flags,omitempty"`
}

// PostAttachment describes an attachment to be uploaded to Bugzilla
type PostAttachment struct {
	ID    int `json:"id,omitempty"`
	BugId int `json:"bug_id,omitempty"`

	Creator        string     `json:"creator,omitempty"`
	Data           []byte     `json:"data,omitempty"`
	IsObsolete     int        `json:"is_obsolete,omitempty"`
	IsPatch        int        `json:"is_patch,omitempty"`
	IsPrivate      int        `json:"is_private,omitempty"`
	CreationTime   *time.Time `json:"creation_time,omitempty"`
	LastChangeTime *time.Time `json:"last_change_time,omitempty"`
	ContentType    string     `json:"content_type,omitempty"`
	Summary        string     `json:"summary,omitempty"`
	Filename       string     `json:"file_name,omitempty"`
	Size           int64      `json:"size,omitempty"`
	Token          string     `json:"token,omitempty"`
	Flags          []Flag     `json:"flags,omitempty"`

	// Only for posting an attachment
	Ids     []int  `json:"ids"`
	Comment string `json:"comment,omitempty"`
}

// Comment as in bug comments
type Comment struct {
	ID           int       `json:"id"`
	BugID        int       `json:"bug_id"`
	AttachmentID *int      `json:"attachment_id,omitempty"`
	Count        int       `json:"count"`
	Text         string    `json:"text"`
	Creator      string    `json:"creator"`
	Time         time.Time `json:"time"`
	CreationTime time.Time `json:"creation_time"`
	IsPrivate    bool      `json:"is_private"`
}

type Fault struct {
	code    int
	message string
}

type responseError struct {
	Code          *int    `json:"code,omitempty"`
	Documentation *string `json:"documentation,omitempty"`
	Error         *bool   `json:"erro,omitempty"`
	Message       *string `json:"message,omitempty"`
}

type bugsResult struct {
	Faults []Fault `json:"faults,omitempty"`
	Bugs   []Bug   `json:"bugs,omitempty"`
}

type commentsResult struct {
	Bugs map[string]map[string][]Comment `json:"bugs"`
}

type bugAttachmentsResult struct {
	Bugs map[string][]Attachment `json:"bugs"`
}

type attachmentsResult struct {
	Attachments map[string]*Attachment `json:"attachments"`
}

func (c *Client) decodeBugs(data []byte) ([]Bug, error) {
	var result bugsResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, DecodeErrror{err}
	}

	return result.Bugs, nil
}

func (c *Client) decodeUpdateResponse(data []byte) ([]UpdateResponse, error) {
	var result updateResponse

	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, ConnectionError{err}
	}

	return result.Bugs, nil
}

func (c *Client) encodeUpdate(update *updBug) ([]byte, error) {
	b, err := json.Marshal(update)
	if err != nil {
		return nil, RequestError{fmt.Errorf("Cannot build update: %v", err)}
	}
	return b, nil
}

func (c *Client) encodePostAttachment(att *PostAttachment) ([]byte, error) {
	b, err := json.Marshal(att)
	if err != nil {
		return nil, RequestError{fmt.Errorf("Cannot build update: %v", err)}
	}
	return b, nil
}

func (c *Client) decodeComments(forBugs []int, data []byte) ([]Comment, error) {
	var result commentsResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, ConnectionError{err}
	}

	comments := make([]Comment, 0, 0)

	// Decode a doc like:
	//
	// const bugsCommentsJson = `
	// {
	//    "bugs" : {
	//       "1047068" : {
	//          "comments" : [
	//             {
	//                "attachment_id" : null,
	//                "bug_id" : 1047068,
	//                "count" : 0,
	//                "creation_time" : "2017-07-03T13:29:15Z",
	//                "creator" : "user1@foobarcorp.example.com",
	//                "id" : 7315202,
	//                "is_private" : false,
	//                "tags" : [],
	//                "text" : "This is a test cloud incident.",
	//                "time" : "2017-07-03T13:29:15Z"
	//             },
	// ...
	for _, forBug := range forBugs {
		strBug := fmt.Sprintf("%d", forBug)
		bugDoc, ok := result.Bugs[strBug]
		if !ok {
			continue
		}
		bugComments, ok := bugDoc["comments"]
		if !ok {
			continue
		}
		for _, comment := range bugComments {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (c *Client) decodeDirectAttachments(data []byte, attachmentIds []int) ([]*Attachment, error) {
	var result attachmentsResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, ConnectionError{err}
	}

	attachments := make([]*Attachment, 0, 0)
	for _, attId := range attachmentIds {
		strAtt := fmt.Sprintf("%d", attId)
		attachment, ok := result.Attachments[strAtt]
		if !ok {
			continue
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func (c *Client) decodeAttachments(data []byte, forBugs []int) ([]Attachment, error) {
	var result bugAttachmentsResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, ConnectionError{err}
	}
	attachments := make([]Attachment, 0, 0)

	for _, forBug := range forBugs {
		strBug := fmt.Sprintf("%d", forBug)
		bugAttachments, ok := result.Bugs[strBug]
		if !ok {
			continue
		}
		for _, attachment := range bugAttachments {
			attachments = append(attachments, attachment)
		}
	}

	return attachments, nil
}

func (c *Client) cacheBugs(bugs []Bug) {
	cacher, ok := c.cacher.(Cacher)
	if !ok {
		return
	}
	for _, bug := range bugs {
		b, err := json.Marshal(bug)
		if err == nil {
			writer := cacher.GetWriter(fmt.Sprintf("%d", bug.ID))
			writer.Write(b)
			writer.Close()
		}
	}
}

func (c *Client) attemptDecodingError(body []byte) (msg string, ok bool) {
	var result responseError
	err := json.Unmarshal(body, &result)
	if err == nil && result.Message != nil {
		msg = fmt.Sprintf("[%d] %s", *result.Code, *result.Message)
		ok = true
	}
	return
}
func (c *Client) collect(resp *http.Response, err error) ([]byte, error) {
	if resp == nil && err != nil {
		return nil, ConnectionError{err}
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	limitedReader := &io.LimitedReader{R: resp.Body, N: 10 * 1024 * 1024}
	body, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return nil, ConnectionError{err}
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		msg, ok := c.attemptDecodingError(body)
		if ok {
			err = fmt.Errorf(msg)
		} else {
			err = fmt.Errorf(http.StatusText(resp.StatusCode))
		}
		return nil, BugzillaError{err}
	}

	return body, nil
}

func (c *Client) fetch(url string) ([]byte, error) {
	resp, err := c.seriousClient.Get(url)
	return c.collect(resp, err)
}

func (c *Client) put(url string, content_type string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.seriousClient.Do(request)
	return c.collect(resp, err)
}

func (c *Client) post(url string, content_type string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.seriousClient.Do(request)
	return c.collect(resp, err)
}

func (c *Client) GetComments(bugIds []int) ([]Comment, error) {
	params := map[string]string{}
	url, err := c.getCommentsURL(bugIds, params)
	if err != nil {
		return nil, err
	}
	body, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	comments, err := c.decodeComments(bugIds, body)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetAttachmentsInfo returns information about attachments in a bug -- with no data
func (c *Client) GetAttachmentsInfo(bugIds []int) ([]Attachment, error) {
	params := map[string]string{"exclude_fields": "data"}
	url, err := c.getAttachmentsURL(bugIds, params)
	if err != nil {
		return nil, err
	}
	body, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	attachments, err := c.decodeAttachments(body, bugIds)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// GetAttachment returns one attachment (including its data)
func (c *Client) GetAttachment(id int) (*Attachment, error) {
	params := map[string]string{}
	url, err := c.getAttachmentURL(id, params)
	if err != nil {
		return nil, err
	}
	body, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	attachments, err := c.decodeDirectAttachments(body, []int{id})
	if err != nil {
		return nil, err
	}
	if len(attachments) != 1 {
		return nil, ConnectionError{fmt.Errorf("Unexpected number of attachments returned: %v", len(attachments))}
	}
	attachment := attachments[0]
	return attachment, nil
}

func (c *Client) GetBug(id int) (*Bug, error) {
	return c.GetBugEx(id, true, true)
}

// GetBug gets a *Bug from the Bugzilla API (apibuzilla)
func (c *Client) GetBugEx(id int, withComments bool, withAttachments bool) (*Bug, error) {
	// query.Set("ctype", "xml")
	// query.Set("excludefield", "attachmentdata")
	params := map[string]string{}
	url, err := c.getBugsURL([]int{id}, params)
	if err != nil {
		return nil, err
	}
	bugsBody, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	bugs, err := c.decodeBugs(bugsBody)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, 0, 0)
	attachments := make([]Attachment, 0, 0)
	if withComments {
		comments, err = c.GetComments([]int{id})
		if err != nil {
			return nil, ConnectionError{err}
		}
	}
	if withAttachments {
		attachments, err = c.GetAttachmentsInfo([]int{id})
		if err != nil {
			return nil, ConnectionError{err}
		}
	}

	for i := range bugs {
		bugs[i].Comments = comments
		bugs[i].Attachments = attachments
	}

	c.cacheBugs(bugs)

	bug := &bugs[0]

	return bug, err
}

// ErrBugzilla is an error from Bugzilla
type ErrBugzilla struct{ error }

func (e ErrBugzilla) Error() string {
	return fmt.Sprintf("Error from Bugzilla: %v", e.error)
}

// PriorityMap maps short priority names to the longer ones, as provided by
// the Web Interface
var PriorityMap = map[string]string{
	"P0": "P0 - Crit Sit",
	"P1": "P1 - Urgent",
	"P2": "P2 - High",
	"P3": "P3 - Medium",
	"P4": "P4 - Low",
	"P5": "P5 - None",
}

// Changes to be performed by Update() for a given bug
type Changes struct {
	SetNeedinfo    string
	RemoveNeedinfo string

	ClearNeedinfo     bool // clears any one needinfo in the bug
	ClearAllNeedinfos bool // + flag ClearNeedinfo to clear all needinfos of the bug
	ClearMyNeedinfos  bool // + flag ClearNeedinfo to clear all my needinfos in the bug

	AddComment       string
	CommentIsPrivate bool

	SetURL         string
	SetAssignee    string
	SetPriority    string
	SetDescription string
	SetWhiteboard  string
	SetStatus      string
	SetResolution  string
	SetDuplicate   int

	AddCc    string
	RemoveCc string
	CcMyself bool

	// DeltaTS should have the timestamp of the last change
	DeltaTS      time.Time
	CheckDeltaTS bool
}

type updateResponse struct {
	Bugs []UpdateResponse `json:"bugs"`
}

// UpdateResponse provides the response from Bugzilla about the changes that
// were submitted
type UpdateResponse struct {
	Alias   []string `json:"alias"`
	Changes map[string]struct {
		Added   string
		Removed string
	} `json:"changes"`
	Id             int       `json:"id"`
	LastChangeTime time.Time `json:"last_change_time"`
}

func (c *Client) checkDeltaTS(changes *Changes, bug *Bug) error {
	if changes.CheckDeltaTS {
		if !bug.LastChangeTime.Equal(changes.DeltaTS) {
			return ErrBugzilla{fmt.Errorf("likely mid-air collision: the bug has been updated at %v", bug.LastChangeTime)}
		}
	}
	return nil
}

type updOp struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
	Set    []string `json:"set,omitempty"`
}

func (u *updOp) AddOp(what string) {
	u.Add = append(u.Add, what)
}

func (u *updOp) RemoveOp(what string) {
	u.Remove = append(u.Remove, what)
}

func (u *updOp) SetOp(what string) {
	u.Set = append(u.Set, what)
}

type flagChange struct {
	Name      string `json:"name,omitempty"`
	TypeID    int    `json:"type_id,omitempty"`
	Status    string `json:"status"`
	Requestee string `json:"requestee,omitempty"`
	ID        int    `json:"id,omitempty"`
	New       bool   `json:"new,omitempty"`
}

type commentChange struct {
	Body      string `json:"body"`
	IsPrivate bool   `json:"is_private"`
}

type updBug struct {
	Ids []int `json:"ids"`

	Alias     []updOp      `json:"alias,omitempty"`
	Blocks    []updOp      `json:"blocks,omitempty"`
	CC        *updOp       `json:"cc,omitempty"`
	CCDetail  []updOp      `json:"cc_detail,omitempty"`
	DependsOn []updOp      `json:"depends_on,omitempty"`
	Flags     []flagChange `json:"flags,omitempty"`
	Groups    []updOp      `json:"groups,omitempty"`
	Keywords  []updOp      `json:"keywords,omitempty"`
	SeeAlso   []updOp      `json:"see_also,omitempty"`

	IsCCAccessible      *bool `json:"is_cc_accessible,omitempty"`
	IsCreatorAccessible *bool `json:"is_creator_accessible,omitempty"`
	ResetAssignedTo     *bool `json:"reset_assigned_to,omitempty"`
	ResetQaContact      *bool `json:"reset_qa_contact,omitempty"`

	AssignedTo      *string        `json:"assigned_to,omitempty"`
	Comment         *commentChange `json:"comment,omitempty"`
	Classification  *string        `json:"classification,omitempty"`
	Component       *string        `json:"component,omitempty"`
	DupeOf          int            `json:"dupe_of,omitempty"`
	OpSys           *string        `json:"op_sys,omitempty"`
	Platform        *string        `json:"platform,omitempty"`
	Priority        *string        `json:"priority,omitempty"`
	Product         *string        `json:"product,omitempty"`
	QAContact       *string        `json:"qa_contact,omitempty"`
	Resolution      *string        `json:"resolution,omitempty"`
	Severity        *string        `json:"severity,omitempty"`
	Status          *string        `json:"status,omitempty"`
	Summary         *string        `json:"summary,omitempty"`
	TargetMilestone *string        `json:"target_milestone,omitempty"`
	URL             *string        `json:"url,omitempty"`
	Version         *string        `json:"version,omitempty"`
	Whiteboard      *string        `json:"whiteboard,omitempty"`

	// Not implemented for now:
	//ActualTime          float64 `json:"actual_time,omitempty"`
	//Deadline            string  `json:"deadline,omitempty"`
	//EstimatedTime       float64 `json:"estimated_time,omitempty"`
	//RemainingTime       float64 `json:"remaining_time,omitempty"`
}

func newBugUpdate() *updBug {
	u := &updBug{}
	u.Ids = make([]int, 0, 0)
	return u
}

func (u *updBug) AddFlagChange(flag *flagChange) {
	u.Flags = append(u.Flags, *flag)
}

func (u *updBug) initUpdOp(target **updOp) {
	if *target == nil {
		*target = &updOp{}
	}
}

func (u *updBug) AddCC(email string) {
	u.initUpdOp(&u.CC)
	u.CC.AddOp(email)
}

func (u *updBug) RemoveCC(email string) {
	u.initUpdOp(&u.CC)
	u.CC.RemoveOp(email)
}

// Update changes a bug with the attribute to be modified provided by
// Changes
func (c *Client) Update(id int, changes Changes) (updateResponse *UpdateResponse, err error) {
	bug, err := c.GetBugEx(id, false, false)
	if err != nil {
		return
	}

	if err = c.checkDeltaTS(&changes, bug); err != nil {
		return
	}

	up := newBugUpdate()
	up.Ids = append(up.Ids, id)

	if changes.AddComment != "" {
		up.Comment = &commentChange{Body: changes.AddComment, IsPrivate: changes.CommentIsPrivate}
	}
	if changes.SetNeedinfo != "" {
		if !bug.hasNeedinfoFor(changes.SetNeedinfo) {
			up.AddFlagChange(&flagChange{
				New:       true,
				Name:      "needinfo",
				Requestee: changes.SetNeedinfo,
				Status:    "?",
			})
		}
	}
	if target := changes.RemoveNeedinfo; target != "" || changes.ClearNeedinfo {
		if changes.ClearMyNeedinfos {
			target, err = c.Config.emailAddress()
			if err != nil {
				return
			}
		} else if changes.ClearNeedinfo {
			target = "" // hint findNeedinfosFor to look for all needinfos
		}
		flagIds := bug.findNeedinfosFor(target)
		if len(flagIds) > 1 && !changes.ClearAllNeedinfos && changes.RemoveNeedinfo == "" {
			err = RequestError{fmt.Errorf("More than one needinfo found")}
			return
		}
		for _, id := range flagIds {
			up.AddFlagChange(&flagChange{
				ID:     id,
				Status: "X",
			})
		}
	}
	if changes.SetURL != "" {
		up.URL = &changes.SetURL
	}
	if changes.SetAssignee != "" {
		up.AssignedTo = &changes.SetAssignee
	}
	if changes.SetDescription != "" {
		up.Summary = &changes.SetDescription
	}
	if changes.SetPriority != "" {
		prio, ok := PriorityMap[changes.SetPriority]
		if !ok {
			err = ErrBugzilla{fmt.Errorf("invalid priority value: %v", changes.SetPriority)}
			return
		}
		up.Priority = &prio
	}
	if changes.AddCc != "" {
		up.AddCC(changes.AddCc)
	}
	if changes.RemoveCc != "" {
		up.RemoveCC(changes.RemoveCc)
	}
	if changes.CcMyself {
		var email string
		email, err = c.Config.emailAddress()
		if err != nil {
			return
		}
		up.AddCC(email)
	}
	if changes.SetWhiteboard != "" {
		up.Whiteboard = &changes.SetWhiteboard
	}
	if changes.SetStatus != "" {
		up.Status = &changes.SetStatus
	}
	if changes.SetResolution != "" {
		up.Resolution = &changes.SetResolution
	}
	if changes.SetDuplicate != 0 {
		up.DupeOf = changes.SetDuplicate
	}

	url, err := c.getUpdateBugsURL([]int{id}, map[string]string{})
	if err != nil {
		return
	}
	updBody, err := c.encodeUpdate(up)
	if err != nil {
		return
	}
	resp, err := c.put(url, "application/json", updBody)
	if err != nil {
		return
	}
	updateResponses, err := c.decodeUpdateResponse(resp)
	if err != nil {
		return
	}

	if len(updateResponses) != 1 {
		return nil, fmt.Errorf("Got an unexpected number of update responses: %v", updateResponses)
	}

	updateResponse = &updateResponses[0]

	return
}

type AttachmentDownload struct {
	client *Client
	id     int

	Attachments map[string]map[string]*Attachment `json:"attachments"`
}

// DecodeFromDownload allows decoding the attachment's data using the raw data
// provided by the download handle from DownloadAttachment
func (a *AttachmentDownload) DataFromDownload(data []byte) ([]byte, error) {
	attachments, err := a.client.decodeDirectAttachments(data, []int{a.id})
	if err != nil {
		return nil, err
	}
	if len(attachments) != 1 {
		return nil, ConnectionError{fmt.Errorf("unexpected number of attachments returned: %v", len(attachments))}
	}
	attachment := attachments[0]
	return attachment.Data, nil
}

// DownloadAttachment an attachment for download
// Returns an Attachment with only the Size and Filename filled, a reader
// and error.
func (c *Client) DownloadAttachment(id int) (*AttachmentDownload, io.ReadCloser, error) {
	params := map[string]string{"includefields": "data"}
	url, err := c.getAttachmentURL(id, params)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.seriousClient.Get(url)
	if err != nil {
		return nil, nil, ConnectionError{err}
	}
	ad := &AttachmentDownload{client: c, id: id}
	return ad, resp.Body, nil
}

// GetBug gets a *Bug from a JSON blob
func (c *Client) GetBugFromJSON(source io.Reader) (*Bug, error) {
	var bug Bug
	decoder := json.NewDecoder(source)
	err := decoder.Decode(&bug)
	if err != nil {
		return nil, err
	}
	return &bug, nil
}

type PostAttachmentResponse struct {
	Ids []int `json:"ids"`
}

func (c *Client) decodePostAttachment(data []byte) (int, error) {
	var result PostAttachmentResponse

	err := json.Unmarshal(data, &result)
	if err != nil {
		return 0, ConnectionError{err}
	}

	return result.Ids[0], nil
}

// UploadAttachment posts a new attachment to a given bug
func (c *Client) UploadAttachment(bugId int, attachment *PostAttachment) (id int, err error) {
	params := map[string]string{}
	url, err := c.getAttachmentsURL([]int{bugId}, params)
	if err != nil {
		return 0, err
	}
	localAttachment := *attachment
	localAttachment.Ids = append(localAttachment.Ids, bugId)

	encoded, err := c.encodePostAttachment(&localAttachment)
	if err != nil {
		return 0, err
	}

	resp, err := c.post(url, "application/json", encoded)
	if err != nil {
		return 0, err
	}
	id, err = c.decodePostAttachment(resp)
	return
}
