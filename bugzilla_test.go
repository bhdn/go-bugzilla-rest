package bugzilla_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kinbiko/jsonassert"

	bugzilla "github.com/bhdn/go-bugzilla-rest"
	. "gopkg.in/check.v1"
)

type clientSuite struct {
	bz *bugzilla.Client
}

var _ = Suite(&clientSuite{})

// Hook up check.v1 into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

func (cs *clientSuite) SetUpTest(c *C) {
}

func makeClient(url string) *bugzilla.Client {
	config := bugzilla.Config{BaseURL: url, ApiKey: "xxxxxx"}
	bz, _ := bugzilla.New(config)
	return bz
}

func (cs *clientSuite) TestCreateClient(c *C) {
	bz := makeClient("http://foobar.com/")
	c.Assert(bz, NotNil)
}

const bugsJson = `
{
   "bugs" : [
      {
         "actual_time" : 0,
         "alias" : [],
         "assigned_to" : "user1@foobarcorp.example.com",
         "assigned_to_detail" : {
            "email" : "user1@foobarcorp.example.com",
            "id" : 63803,
            "name" : "user1@foobarcorp.example.com",
            "real_name" : "Firstname1 LastName1"
         },
         "blocks" : [],
         "cc" : [
            "561726581864@foobarcorp.example.com",
            "user2@foobarcorp.example.com",
            "user1@foobarcorp.example.com",
            "user3@foobarcorp.example.com"
         ],
         "cc_detail" : [
            {
               "email" : "user2@foobarcorp.example.com",
               "id" : 67231,
               "name" : "user2@foobarcorp.example.com",
               "real_name" : "Firstname2 Lastname2"
            },
            {
               "email" : "user1@foobarcorp.example.com",
               "id" : 63803,
               "name" : "user1@foobarcorp.example.com",
               "real_name" : "Firstname1 LastName1"
            },
            {
               "email" : "user3@foobarcorp.example.com",
               "id" : 91277,
               "name" : "user3@foobarcorp.example.com",
               "real_name" : "Firstname3 Lastname3"
            }
         ],
         "cf_biz_priority" : "",
         "cf_blocker" : "---",
         "cf_foundby" : "i18n Test",
         "cf_it_deployment" : "---",
         "cf_marketing_qa_status" : "---",
         "cf_nts_priority" : "",
         "classification" : "Enterprise Frobnicator",
         "component" : "Basesystem",
         "creation_time" : "2017-07-03T13:29:15Z",
         "creator" : "user1@foobarcorp.example.com",
         "creator_detail" : {
            "email" : "user1@foobarcorp.example.com",
            "id" : 63803,
            "name" : "user1@foobarcorp.example.com",
            "real_name" : "Firstname1 LastName1"
         },
         "deadline" : null,
         "depends_on" : [],
         "dupe_of" : null,
         "estimated_time" : 0,
         "flags" : [
            {
               "creation_date" : "2022-04-26T07:53:49Z",
               "id" : 264343,
               "modification_date" : "2022-04-26T07:53:49Z",
               "name" : "needinfo",
               "requestee" : "user1@foobarcorp.example.com",
               "setter" : "user1@foobarcorp.example.com",
               "status" : "?",
               "type_id" : 4
            },
            {
               "creation_date" : "2022-06-14T15:54:28Z",
               "id" : 266294,
               "modification_date" : "2022-06-14T15:54:28Z",
               "name" : "needinfo",
               "requestee" : "user3@foobarcorp.example.com",
               "setter" : "user3@foobarcorp.example.com",
               "status" : "?",
               "type_id" : 4
            },
            {
               "creation_date" : "2022-06-14T15:54:28Z",
               "id" : 266299,
               "modification_date" : "2022-06-14T15:54:28Z",
               "name" : "needinfo",
               "requestee" : "user3@foobarcorp.example.com",
               "setter" : "user3@foobarcorp.example.com",
               "status" : "?",
               "type_id" : 4
            },
            {
               "creation_date" : "2019-03-27T13:50:29Z",
               "id" : 201663,
               "modification_date" : "2019-03-27T13:50:29Z",
               "name" : "SHIP_STOPPER",
               "requestee" : "user1@foobarcorp.example.com",
               "setter" : "user1@foobarcorp.example.com",
               "status" : "?",
               "type_id" : 2
            },
            {
               "creation_date" : "2022-04-26T08:05:37Z",
               "id" : 264345,
               "modification_date" : "2022-04-26T08:05:37Z",
               "name" : "CCB_Review",
               "setter" : "user2@foobarcorp.example.com",
               "status" : "+",
               "type_id" : 3
            }
         ],
         "groups" : [
            "foobarcorponly",
            "FOOBARCorp Enterprise Partner"
         ],
         "id" : 1047068,
         "is_cc_accessible" : true,
         "is_confirmed" : true,
         "is_creator_accessible" : true,
         "is_open" : true,
         "keywords" : [
            "TRETA",
            "TRETA_ADDRESSED"
         ],
         "last_change_time" : "2023-04-12T01:02:03Z",
         "op_sys" : "Other",
         "platform" : "Other",
         "priority" : "P2 - High",
         "product" : "Enterprise Frobnicator 9000.1",
         "qa_contact" : "user1@foobarcorp.example.com",
         "qa_contact_detail" : {
            "email" : "user1@foobarcorp.example.com",
            "id" : 63803,
            "name" : "user1@foobarcorp.example.com",
            "real_name" : "Firstname1 LastName1"
         },
         "remaining_time" : 0,
         "resolution" : "",
         "see_also" : [],
         "severity" : "Major",
         "status" : "REOPENED",
         "summary" : "L4: test cloud bug123",
         "target_milestone" : "FROB90001Maint-Upd",
         "update_token" : "1683306765-PMQ3v1SB5rHQwTPnDeSPrCAmChAk5itzZn7A_WfGgq4",
         "url" : "https://xxxxxx.foobarcorp.example.com/incident/9999999",
         "version" : "FROB90001Maint-Upd",
         "whiteboard" : "wasXXXFLAG:48626 zzz wasXXXFLAG:54027 é com acento wasXXXFLAG:57166 openXXXFLAG:59072 wasXXXFLAG:59407 wasXXXFLAG:59403 wasXXXFLAG:59411 wasXXXFLAG:63174 wasXXXFLAG:63342 wasXXXFLAG:65140 wasXXXFLAG:65138 wasXXXFLAG:65421 QE_REVIEW wasXXXFLAG:65669 l3bs:c1,c2,c9"
      }
   ],
   "faults" : []
}
`

const bugsCommentsJson = `
{
   "bugs" : {
      "1047068" : {
         "comments" : [
            {
               "attachment_id" : null,
               "bug_id" : 1047068,
               "count" : 0,
               "creation_time" : "2017-07-03T13:29:15Z",
               "creator" : "user1@foobarcorp.example.com",
               "id" : 7315202,
               "is_private" : false,
               "tags" : [],
               "text" : "This is a test cloud incident.",
               "time" : "2017-07-03T13:29:15Z"
            },
            {
               "attachment_id" : null,
               "bug_id" : 1047068,
               "count" : 1,
               "creation_time" : "2017-07-03T13:31:23Z",
               "creator" : "bot1@foobarcorp.example.com",
               "id" : 7315205,
               "is_private" : true,
               "tags" : [],
               "text" : "XXXFLAG:48626 is now handled by Firstname1 LastName1.",
               "time" : "2017-07-03T13:31:23Z"
            },
            {
               "attachment_id" : null,
               "bug_id" : 1047068,
               "count" : 2,
               "creation_time" : "2017-07-11T10:50:17Z",
               "creator" : "user1@foobarcorp.example.com",
               "id" : 7323867,
               "is_private" : false,
               "tags" : [],
               "text" : "This is a multi line comment.\n\nThe purpose is to test how far it can go with multiline handling. Word Word Word Word Word Word Word Word Word Word Word Word Word Word. \n\nDone.",
               "time" : "2017-07-11T10:50:17Z"
            },
            {
               "attachment_id" : null,
               "bug_id" : 1047068,
               "count" : 3,
               "creation_time" : "2017-07-11T11:00:00Z",
               "creator" : "user1@foobarcorp.example.com",
               "id" : 7323889,
               "is_private" : false,
               "tags" : [],
               "text" : "This is a multi line comment.\n\nThe purpose is to test how far it can go with multiline handling. Word Word Word Word Word Word Word Word Word Word Word Word Word Word. \n\nDone.",
               "time" : "2017-07-11T11:00:00Z"
            }
         ]
      }
   },
   "comments" : {}
}
`

var bugsAttachmentsJson = `
{
   "attachments" : {},
   "bugs" : {
      "1047068" : [
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T12:48:24Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766283,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 0,
            "last_change_time" : "2018-04-06T12:48:24Z",
            "size" : 2,
            "summary" : "description"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T12:50:44Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766284,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 0,
            "last_change_time" : "2018-04-06T12:50:44Z",
            "size" : 2,
            "summary" : "description"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T12:58:52Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766285,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 0,
            "last_change_time" : "2018-04-06T12:58:52Z",
            "size" : 2,
            "summary" : "description"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T13:02:41Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766286,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 0,
            "last_change_time" : "2018-04-06T13:02:41Z",
            "size" : 2,
            "summary" : "description3"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T13:07:27Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766287,
            "is_obsolete" : 0,
            "is_patch" : 1,
            "is_private" : 0,
            "last_change_time" : "2018-04-06T13:07:27Z",
            "size" : 2,
            "summary" : "description3"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T13:16:31Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766288,
            "is_obsolete" : 0,
            "is_patch" : 1,
            "is_private" : 1,
            "last_change_time" : "2018-04-06T13:16:31Z",
            "size" : 2,
            "summary" : "description4"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T13:16:50Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766289,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 1,
            "last_change_time" : "2018-04-06T13:16:50Z",
            "size" : 2,
            "summary" : "description5test"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2018-04-06T13:24:10Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "a.txt",
            "flags" : [],
            "id" : 766292,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 1,
            "last_change_time" : "2018-04-06T13:24:10Z",
            "size" : 2,
            "summary" : "description5test"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2019-11-12T17:46:52Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "foo",
            "flags" : [],
            "id" : 823972,
            "is_obsolete" : 0,
            "is_patch" : 1,
            "is_private" : 0,
            "last_change_time" : "2019-11-12T17:46:52Z",
            "size" : 533,
            "summary" : "Another test attachment é"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2019-11-12T17:49:55Z",
            "creator" : "user1@foobarcorp.example.com",
            "file_name" : "foo",
            "flags" : [],
            "id" : 823974,
            "is_obsolete" : 0,
            "is_patch" : 1,
            "is_private" : 0,
            "last_change_time" : "2019-11-12T17:49:55Z",
            "size" : 533,
            "summary" : "Another test attachment é com acento"
         },
         {
            "bug_id" : 1047068,
            "content_type" : "text/plain",
            "creation_time" : "2022-05-31T20:04:05Z",
            "creator" : "user2@foobarcorp.example.com",
            "file_name" : "test.txt",
            "flags" : [],
            "id" : 859321,
            "is_obsolete" : 0,
            "is_patch" : 0,
            "is_private" : 0,
            "last_change_time" : "2022-05-31T20:04:05Z",
            "size" : 5,
            "summary" : "Test attachment"
         }
      ]
   }
}
`

const singleAttachmentJson = `
{
   "attachments" : {
      "766288" : {
         "data" : "YQo="
      }
   },
   "bugs" : {}
}
`

func (cs *clientSuite) makeBugzillaRestServer(c *C, bugId int) *httptest.Server {
	bncPart := fmt.Sprintf("/%d", bugId)
	bncPath := "/rest/bug" + bncPart
	comPath := "/rest/bug" + bncPart + "/comment"
	attPath := "/rest/bug" + bncPart + "/attachment"
	attOnlyPath := "/rest/bug/attachment/766288"
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case bncPath:
			c.Assert(strings.HasSuffix(r.URL.EscapedPath(), bncPart), Equals, true)
			io.WriteString(w, bugsJson)
		case comPath:
			//c.Assert(strings.Contains(r.URL.EscapedPath(), bncPart), Equals, true)
			io.WriteString(w, bugsCommentsJson)
		case attPath:
			//c.Assert(strings.Contains(r.URL.EscapedPath(), bncPart), Equals, true)
			io.WriteString(w, bugsAttachmentsJson)
		case attOnlyPath:
			io.WriteString(w, singleAttachmentJson)
		default:
			msg := fmt.Sprintf("Unimplemented for %s", r.URL.Path)
			http.Error(w, msg, 500)
			return
		}
	}))
	return ts0
}

func (cs *clientSuite) TestGetBug(c *C) {
	ts0 := cs.makeBugzillaRestServer(c, 1047068)
	defer ts0.Close()

	bz := makeClient(ts0.URL)
	bug, err := bz.GetBug(1047068)
	c.Assert(err, IsNil)
	c.Assert(bug, NotNil)
	c.Assert(bug.ID, Equals, 1047068)
	c.Assert(bug.Summary, Equals, "L4: test cloud bug123")
	c.Assert(bug.CreationTime, Equals, time.Date(2017, 7, 3, 13, 29, 15, 0, time.UTC))
	c.Assert(bug.LastChangeTime, Equals, time.Date(2023, 4, 12, 01, 02, 03, 0, time.UTC))
	c.Assert(bug.IsCCAccessible, Equals, true)
	c.Assert(bug.IsConfirmed, Equals, true)
	c.Assert(bug.IsCreatorAccessible, Equals, true)
	c.Assert(bug.IsOpen, Equals, true)
	c.Assert(bug.Classification, Equals, "Enterprise Frobnicator")
	c.Assert(bug.Component, Equals, "Basesystem")
	c.Assert(bug.Version, Equals, "FROB90001Maint-Upd")
	c.Assert(bug.Platform, Equals, "Other")
	c.Assert(bug.OpSys, Equals, "Other")
	c.Assert(bug.Priority, Equals, "P2 - High")
	c.Assert(bug.Product, Equals, "Enterprise Frobnicator 9000.1")
	c.Assert(bug.QAContact, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.RemainingTime, Equals, 0.0)
	c.Assert(bug.Resolution, Equals, "")
	c.Assert(bug.Severity, Equals, "Major")
	c.Assert(bug.Status, Equals, "REOPENED")
	c.Assert(bug.TargetMilestone, Equals, "FROB90001Maint-Upd")
	c.Assert(bug.URL, Equals, "https://xxxxxx.foobarcorp.example.com/incident/9999999")
	c.Assert(bug.UpdateToken, Equals, "1683306765-PMQ3v1SB5rHQwTPnDeSPrCAmChAk5itzZn7A_WfGgq4")
	c.Assert(bug.Whiteboard, Equals, "wasXXXFLAG:48626 zzz wasXXXFLAG:54027 é com acento wasXXXFLAG:57166 openXXXFLAG:59072 wasXXXFLAG:59407 wasXXXFLAG:59403 wasXXXFLAG:59411 wasXXXFLAG:63174 wasXXXFLAG:63342 wasXXXFLAG:65140 wasXXXFLAG:65138 wasXXXFLAG:65421 QE_REVIEW wasXXXFLAG:65669 l3bs:c1,c2,c9")

	c.Assert(len(bug.Comments), Equals, 4)
	c.Assert(len(bug.Comments), Equals, 4)
	c.Assert(bug.Comments[0].ID, Equals, 7315202)
	c.Assert(bug.Comments[0].BugID, Equals, 1047068)
	c.Assert(bug.Comments[0].AttachmentID, IsNil)
	c.Assert(bug.Comments[0].Count, Equals, 0)
	c.Assert(bug.Comments[0].Text, Equals, "This is a test cloud incident.")
	c.Assert(bug.Comments[0].Creator, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.Comments[0].Time, Equals, time.Date(2017, 7, 3, 13, 29, 15, 0, time.UTC))
	c.Assert(bug.Comments[0].CreationTime, Equals, time.Date(2017, 7, 3, 13, 29, 15, 0, time.UTC))
	c.Assert(bug.Comments[0].IsPrivate, Equals, false)

	c.Assert(bug.Comments[1].ID, Equals, 7315205)
	c.Assert(bug.Comments[1].BugID, Equals, 1047068)
	c.Assert(bug.Comments[1].AttachmentID, IsNil)
	c.Assert(bug.Comments[1].Count, Equals, 1)
	c.Assert(bug.Comments[1].Text, Equals, "XXXFLAG:48626 is now handled by Firstname1 LastName1.")
	c.Assert(bug.Comments[1].Creator, Equals, "bot1@foobarcorp.example.com")
	c.Assert(bug.Comments[1].Time, Equals, time.Date(2017, 7, 3, 13, 31, 23, 0, time.UTC))
	c.Assert(bug.Comments[1].CreationTime, Equals, time.Date(2017, 7, 3, 13, 31, 23, 0, time.UTC))
	c.Assert(bug.Comments[1].IsPrivate, Equals, true)

	c.Assert(bug.Comments[2].ID, Equals, 7323867)
	c.Assert(bug.Comments[2].BugID, Equals, 1047068)
	c.Assert(bug.Comments[2].AttachmentID, IsNil)
	c.Assert(bug.Comments[2].Count, Equals, 2)
	c.Assert(bug.Comments[2].Text, Equals, "This is a multi line comment.\n\nThe purpose is to test how far it can go with multiline handling. Word Word Word Word Word Word Word Word Word Word Word Word Word Word. \n\nDone.")
	c.Assert(bug.Comments[2].Creator, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.Comments[2].Time, Equals, time.Date(2017, 7, 11, 10, 50, 17, 0, time.UTC))
	c.Assert(bug.Comments[2].CreationTime, Equals, time.Date(2017, 7, 11, 10, 50, 17, 0, time.UTC))
	c.Assert(bug.Comments[2].IsPrivate, Equals, false)

	c.Assert(bug.Comments[3].ID, Equals, 7323889)
	c.Assert(bug.Comments[3].BugID, Equals, 1047068)
	c.Assert(bug.Comments[3].AttachmentID, IsNil)
	c.Assert(bug.Comments[3].Count, Equals, 3)
	c.Assert(bug.Comments[3].Text, Equals, "This is a multi line comment.\n\nThe purpose is to test how far it can go with multiline handling. Word Word Word Word Word Word Word Word Word Word Word Word Word Word. \n\nDone.")

	// Flags
	c.Assert(len(bug.Flags), Equals, 5)

	// Check the properties of each flag
	for _, f := range bug.Flags {
		switch f.ID {
		case 264343:
			c.Assert(f.Name, Equals, "needinfo")
			c.Assert(f.TypeID, Equals, 4)
			c.Assert(f.CreationDate, Equals, time.Date(2022, 4, 26, 7, 53, 49, 0, time.UTC))
			c.Assert(f.ModificationDate, Equals, time.Date(2022, 4, 26, 7, 53, 49, 0, time.UTC))
			c.Assert(f.Status, Equals, "?")
			c.Assert(f.Setter, Equals, "user1@foobarcorp.example.com")
			c.Assert(f.Requestee, Equals, "user1@foobarcorp.example.com")
		case 266294:
			c.Assert(f.Name, Equals, "needinfo")
			c.Assert(f.TypeID, Equals, 4)
			c.Assert(f.CreationDate, Equals, time.Date(2022, 6, 14, 15, 54, 28, 0, time.UTC))
			c.Assert(f.ModificationDate, Equals, time.Date(2022, 6, 14, 15, 54, 28, 0, time.UTC))
			c.Assert(f.Status, Equals, "?")
			c.Assert(f.Setter, Equals, "user3@foobarcorp.example.com")
			c.Assert(f.Requestee, Equals, "user3@foobarcorp.example.com")
		case 201663:
			c.Assert(f.Name, Equals, "SHIP_STOPPER")
			c.Assert(f.TypeID, Equals, 2)
			c.Assert(f.CreationDate, Equals, time.Date(2019, 3, 27, 13, 50, 29, 0, time.UTC))
			c.Assert(f.ModificationDate, Equals, time.Date(2019, 3, 27, 13, 50, 29, 0, time.UTC))
			c.Assert(f.Status, Equals, "?")
			c.Assert(f.Setter, Equals, "user1@foobarcorp.example.com")
			c.Assert(f.Requestee, Equals, "user1@foobarcorp.example.com")
		case 264345:
			c.Assert(f.Name, Equals, "CCB_Review")
			c.Assert(f.TypeID, Equals, 3)
			c.Assert(f.CreationDate, Equals, time.Date(2022, 4, 26, 8, 5, 37, 0, time.UTC))
			c.Assert(f.ModificationDate, Equals, time.Date(2022, 4, 26, 8, 5, 37, 0, time.UTC))
			c.Assert(f.Status, Equals, "+")
			c.Assert(f.Setter, Equals, "user2@foobarcorp.example.com")
			c.Assert(f.Requestee, Equals, "")
		}
	}

	// Attachments
	c.Assert(len(bug.Attachments), Equals, 11)
	c.Assert(bug.Attachments[0].ID, Equals, 766283)
	c.Assert(bug.Attachments[0].BugId, Equals, 1047068)
	c.Assert(bug.Attachments[0].Creator, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.Attachments[0].IsObsolete, Equals, 0)
	c.Assert(bug.Attachments[0].IsPatch, Equals, 0)
	c.Assert(bug.Attachments[0].IsPrivate, Equals, 0)
	c.Assert(bug.Attachments[0].CreationTime, Equals, time.Date(2018, 4, 6, 12, 48, 24, 0, time.UTC))
	c.Assert(bug.Attachments[0].LastChangeTime, Equals, time.Date(2018, 4, 6, 12, 48, 24, 0, time.UTC))
	c.Assert(bug.Attachments[0].ContentType, Equals, "text/plain")
	c.Assert(bug.Attachments[0].Summary, Equals, "description")
	c.Assert(bug.Attachments[0].Filename, Equals, "a.txt")
	c.Assert(bug.Attachments[0].Size, Equals, int64(2))
	c.Assert(bug.Attachments[0].Token, Equals, "")
	c.Assert(bug.Attachments[0].Flags, HasLen, 0)
	c.Assert(bug.Attachments[1].ID, Equals, 766284)
	c.Assert(bug.Attachments[1].BugId, Equals, 1047068)
	c.Assert(bug.Attachments[1].Creator, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.Attachments[1].IsObsolete, Equals, 0)
	c.Assert(bug.Attachments[1].IsPatch, Equals, 0)
	c.Assert(bug.Attachments[1].IsPrivate, Equals, 0)
	c.Assert(bug.Attachments[1].CreationTime, Equals, time.Date(2018, 4, 6, 12, 50, 44, 0, time.UTC))
	c.Assert(bug.Attachments[1].LastChangeTime, Equals, time.Date(2018, 4, 6, 12, 50, 44, 0, time.UTC))
	c.Assert(bug.Attachments[1].ContentType, Equals, "text/plain")
	c.Assert(bug.Attachments[1].Summary, Equals, "description")
	c.Assert(bug.Attachments[1].Filename, Equals, "a.txt")
	c.Assert(bug.Attachments[1].Size, Equals, int64(2))
	c.Assert(bug.Attachments[1].Token, Equals, "")
	c.Assert(bug.Attachments[1].Flags, HasLen, 0)
	c.Assert(bug.Attachments[2].BugId, Equals, 1047068)
	c.Assert(bug.Attachments[2].ContentType, Equals, "text/plain")
	c.Assert(bug.Attachments[2].CreationTime, Equals, time.Date(2018, 4, 6, 12, 58, 52, 0, time.UTC))
	c.Assert(bug.Attachments[2].Creator, Equals, "user1@foobarcorp.example.com")
	c.Assert(bug.Attachments[2].Filename, Equals, "a.txt")
	c.Assert(bug.Attachments[2].Flags, HasLen, 0)
	c.Assert(bug.Attachments[2].ID, Equals, 766285)
	c.Assert(bug.Attachments[2].IsObsolete, Equals, 0)
	c.Assert(bug.Attachments[2].IsPatch, Equals, 0)
	c.Assert(bug.Attachments[2].IsPrivate, Equals, 0)
	c.Assert(bug.Attachments[2].LastChangeTime, Equals, time.Date(2018, 4, 6, 12, 58, 52, 0, time.UTC))
	c.Assert(bug.Attachments[2].Size, Equals, int64(2))
	c.Assert(bug.Attachments[2].Summary, Equals, "description")
}

func makeClientWithCache(url string, cacher bugzilla.Cacher) *bugzilla.Client {
	config := bugzilla.Config{BaseURL: url, ApiKey: "xxxxxxxxxxx", Cacher: cacher}
	bz, _ := bugzilla.New(config)
	return bz
}

type CacherHelper struct {
	bugzilla.Cacher

	location string
	buf      FakeBuf
	id       string
}

type FakeBuf struct {
	bytes.Buffer
}

func (f *FakeBuf) Close() error {
	return nil
}

func (c *CacherHelper) GetWriter(id string) io.WriteCloser {
	c.id = id
	return &c.buf
}

type cacheChecker struct {
	ID               int            `json:"id"`
	Summary          string         `json:"summary"`
	AssignedToDetail *bugzilla.User `json:"assigned_to_detail"`
}

func (cs *clientSuite) TestGetBugWithCacher(c *C) {
	ts0 := cs.makeBugzillaRestServer(c, 1047068)
	defer ts0.Close()

	var cacher CacherHelper
	bz := makeClientWithCache(ts0.URL, &cacher)
	bug, err := bz.GetBug(1047068)
	c.Assert(err, IsNil)
	c.Assert(cacher.id, Equals, "1047068")
	c.Assert(bug.ID, Equals, 1047068)
	var checker cacheChecker
	err = json.Unmarshal(cacher.buf.Bytes(), &checker)
	c.Assert(err, IsNil)
	c.Check(checker.ID, Equals, 1047068)
	c.Check(checker.Summary, Equals, "L4: test cloud bug123")
	c.Check(checker.AssignedToDetail.RealName, Equals, "Firstname1 LastName1")
}

func makeBugzillaServerWithChannels() (*httptest.Server, chan string, chan string, chan string) {
	queries := make(chan string, 10)
	nextJson := make(chan string, 10)
	processBug := make(chan string, 10)
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			io.WriteString(w, <-nextJson)
			return
		} else if r.Method == "PUT" || r.Method == "POST" {
			buf := new(strings.Builder)
			io.Copy(buf, r.Body)
			queries <- buf.String()
			io.WriteString(w, <-processBug)
			return
		}
		http.Error(w, "Unimplemented", 500)
	}))
	return ts0, queries, nextJson, processBug
}

func (cs *clientSuite) TestUpdateBug(c *C) {
	ts0, queries, nextJson, processBug := makeBugzillaServerWithChannels()
	defer ts0.Close()
	bz := makeClient(ts0.URL)

	email := "user@foobar.com"
	changes := bugzilla.Changes{SetNeedinfo: email}
	nextJson <- bugsJson
	processBug <- ` { "bugs" : [ {
         "alias" : [],
         "changes" : {
            "flagtypes.name" : {
               "added" : "",
               "removed" : "needinfo?(user@foobar.com)"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-09T09:30:30Z"
      } ] }
	`
	result, err := bz.Update(101234, changes)
	c.Assert(err, IsNil)
	ja := jsonassert.New(c)
	query := <-queries
	ja.Assertf(query, `{"ids": [101234], "flags": [{"name": "needinfo", "new": true, "requestee": "user@foobar.com", "status": "?"}]}`)
	c.Assert(result.Id, Equals, 101234)
	c.Assert(len(result.Alias), Equals, 0)
	c.Assert(len(result.Changes), Equals, 1)

	// Ensure no changes in Changes{} means no changes in the bug
	changes = bugzilla.Changes{}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {},
         "id" : 101234,
         "last_change_time" : "2023-05-09T09:53:05Z"
      }
   ]
}	`
	result, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234]}`)

	comment := "this is a comment"
	changes = bugzilla.Changes{AddComment: comment}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {},
         "id" : 101234,
         "last_change_time" : "2023-05-09T09:53:05Z"
      }
   ]
}	`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "comment": {"body": "this is a comment", "is_private": false}}`)

	// Ensure the comment can be marked as private
	changes = bugzilla.Changes{AddComment: comment, CommentIsPrivate: true}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {},
         "id" : 101234,
         "last_change_time" : "2023-05-09T09:53:05Z"
      }
   ]
}	`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "comment": {"body": "this is a comment", "is_private": true}}`)

	// RemoveNeedinfo clears all needinfos of a given email address
	changes = bugzilla.Changes{RemoveNeedinfo: "user3@foobarcorp.example.com"}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {
            "flagtypes.name" : {
               "added" : "",
               "removed" : "needinfo?(user3@foobarcorp.example.com), needinfo?(user3@foobarcorp.example.com)"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-09T17:38:49Z"
      }
   ]
}`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "flags": [{"status": "X", "id": 266294}, {"status": "X", "id": 266299}]}`)

	// ClearNeedinfo fails if more than one needinfo is set and ClearAllNeedinfos is not
	changes = bugzilla.Changes{ClearNeedinfo: true}
	nextJson <- bugsJson
	_, err = bz.Update(101234, changes)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, ".*More than one needinfo.*")

	// Ensure ClearAllNeedinfos filters all needinfos
	changes = bugzilla.Changes{ClearNeedinfo: true, ClearAllNeedinfos: true}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {
            "flagtypes.name" : {
               "added" : "",
               "removed" : "needinfo?(user1@foobarcorp.example.com), needinfo?(user3@foobarcorp.example.com), needinfo?(user3@foobarcorp.example.com)"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-09T17:38:49Z"
      }
   ]
}`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "flags": [{"status": "X", "id": 264343}, {"status": "X", "id": 266294}, {"status": "X", "id": 266299}]}`)

	// ClearNeedinfo uses the Email or Username in Config if ClearMyNeedinfos is set
	oldConfig := bz.Config
	bz.Config.Username = "user1@foobarcorp.example.com"
	changes = bugzilla.Changes{ClearNeedinfo: true, ClearMyNeedinfos: true}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {
            "flagtypes.name" : {
               "added" : "",
               "removed" : "needinfo?(user1@foobarcorp.example.com)"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-09T17:38:49Z"
      }
   ]
}`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "flags": [{"status": "X", "id": 264343}]}`)
	bz.Config = oldConfig

	url := "http://foobar.com/1/2"
	changes = bugzilla.Changes{SetURL: url}
	nextJson <- bugsJson
	processBug <- `
{
   "bugs" : [
      {
         "alias" : [],
         "changes" : {
            "url" : {
               "added" : "https://sample.example.com/foo/bar/baz",
               "removed" : "https://l3support.foobarcorp.example.com/incident/59072"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-10T09:32:26Z"
      }
   ]
}
`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "url": "http://foobar.com/1/2"}`)

	priority := "P0"
	changes = bugzilla.Changes{SetPriority: priority}
	nextJson <- bugsJson
	processBug <- ` { "bugs" : [ {
         "alias" : [],
         "changes" : {
            "priority" : {
               "added" : "P0 - Crit Sit",
               "removed" : "P2 - High"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-10T09:30:25Z"
      } ] } `
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "priority": "P0 - Crit Sit"}`)

	priority = "wrong"
	changes = bugzilla.Changes{SetPriority: priority}
	nextJson <- bugsJson
	_, err = bz.Update(101234, changes)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, ".*invalid priority value.*")

	email = "user4@foobarcorp.example.com"
	changes = bugzilla.Changes{AddCc: email}
	nextJson <- bugsJson
	processBug <- `{"bugs" : [ {
         "alias" : [],
         "changes" : {
            "cc" : {
               "added" : "user4@foobarcorp.example.com",
               "removed" : ""
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-10T10:00:02Z"
      } ] }`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "cc": {"add": ["user4@foobarcorp.example.com"]}}`)

	email = "user4@foobarcorp.example.com"
	changes = bugzilla.Changes{RemoveCc: email}
	nextJson <- bugsJson
	processBug <- `{ "bugs" : [ { "alias" : [],
         "changes" : {
            "cc" : {
               "added" : "",
               "removed" : "debogdano@gmail.com"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-10T10:02:47Z"
      } ] }`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "cc": {"remove": ["user4@foobarcorp.example.com"]}}`)

	oldConfig = bz.Config
	bz.Config.Username = "user5@foobarcorp.example.com"
	changes = bugzilla.Changes{CcMyself: true}
	nextJson <- bugsJson
	processBug <- `{ "bugs" : [ { "alias" : [],
         "changes" : {
            "cc" : {
               "added" : "user5@foobarcorp.example.com",
               "removed" : "debogdano@gmail.com"
            }
         },
         "id" : 101234,
         "last_change_time" : "2023-05-10T10:02:47Z"
      } ] }`
	_, err = bz.Update(101234, changes)
	c.Assert(err, IsNil)
	query = <-queries
	ja.Assertf(query, `{"ids": [101234], "cc": {"add": ["user5@foobarcorp.example.com"]}}`)
	bz.Config = oldConfig

	// Fail when Username is not an email
	oldConfig = bz.Config
	bz.Config.Username = "user5"
	changes = bugzilla.Changes{CcMyself: true}
	nextJson <- bugsJson
	_, err = bz.Update(101234, changes)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, ".*Your username doesn't look like an email address.*")
	bz.Config = oldConfig

	lastChange := time.Date(2023, 04, 12, 02, 02, 03, 0, time.UTC)
	changes = bugzilla.Changes{AddComment: "Some comment", DeltaTS: lastChange, CheckDeltaTS: true}
	nextJson <- bugsJson
	_, err = bz.Update(101234, changes)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, ".*collision.*")

}

func (cs *clientSuite) TestUnauthorized(c *C) {
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer ts0.Close()
	bz := makeClient(ts0.URL)
	bug, err := bz.GetBug(1047068)
	c.Assert(bug, IsNil)
	c.Assert(err, ErrorMatches, ".*Unauthorized*")
}

var sampleError = `
{
   "code" : 102,
   "documentation" : "https://bugzilla.readthedocs.org/en/5.0/api/",
   "error" : true,
   "message" : "You are not authorized to access bug #1171184."
}
`

// TestGetBugNotPermitted triggers the code path when a request is made to
// the wrong endpoint of Bugzilla (one that uses another type of auth)
func (cs *clientSuite) TestGetBugNotPermitted(c *C) {
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, sampleError, http.StatusBadRequest)
	}))
	defer ts0.Close()
	bz := makeClient(ts0.URL)
	bug, err := bz.GetBug(1047068)
	c.Assert(bug, IsNil)
	c.Assert(err, ErrorMatches, ".*You are not authorized to access bug.*")
}

// A copy of bugsJson, modified to have only one bug document
var sampleJSON = `
{
 "actual_time" : 0,
 "alias" : [],
 "assigned_to" : "user1@foobarcorp.example.com",
 "assigned_to_detail" : {
    "email" : "user1@foobarcorp.example.com",
    "id" : 63803,
    "name" : "user1@foobarcorp.example.com",
    "real_name" : "Firstname1 LastName1"
 },
 "blocks" : [],
 "cc" : [
    "561726581864@foobarcorp.example.com",
    "user2@foobarcorp.example.com",
    "user1@foobarcorp.example.com",
    "user3@foobarcorp.example.com"
 ],
 "cc_detail" : [
    {
       "email" : "user2@foobarcorp.example.com",
       "id" : 67231,
       "name" : "user2@foobarcorp.example.com",
       "real_name" : "Firstname2 Lastname2"
    },
    {
       "email" : "user1@foobarcorp.example.com",
       "id" : 63803,
       "name" : "user1@foobarcorp.example.com",
       "real_name" : "Firstname1 LastName1"
    },
    {
       "email" : "user3@foobarcorp.example.com",
       "id" : 91277,
       "name" : "user3@foobarcorp.example.com",
       "real_name" : "Firstname3 Lastname3"
    }
 ],
 "cf_biz_priority" : "",
 "cf_blocker" : "---",
 "cf_foundby" : "i18n Test",
 "cf_it_deployment" : "---",
 "cf_marketing_qa_status" : "---",
 "cf_nts_priority" : "",
 "classification" : "Enterprise Frobnicator",
 "component" : "Basesystem",
 "creation_time" : "2017-07-03T13:29:15Z",
 "creator" : "user1@foobarcorp.example.com",
 "creator_detail" : {
    "email" : "user1@foobarcorp.example.com",
    "id" : 63803,
    "name" : "user1@foobarcorp.example.com",
    "real_name" : "Firstname1 LastName1"
 },
 "deadline" : null,
 "depends_on" : [],
 "dupe_of" : null,
 "estimated_time" : 0,
 "flags" : [
    {
       "creation_date" : "2022-04-26T07:53:49Z",
       "id" : 264343,
       "modification_date" : "2022-04-26T07:53:49Z",
       "name" : "needinfo",
       "requestee" : "user1@foobarcorp.example.com",
       "setter" : "user1@foobarcorp.example.com",
       "status" : "?",
       "type_id" : 4
    },
    {
       "creation_date" : "2022-06-14T15:54:28Z",
       "id" : 266294,
       "modification_date" : "2022-06-14T15:54:28Z",
       "name" : "needinfo",
       "requestee" : "user3@foobarcorp.example.com",
       "setter" : "user3@foobarcorp.example.com",
       "status" : "?",
       "type_id" : 4
    },
    {
       "creation_date" : "2022-06-14T15:54:28Z",
       "id" : 266299,
       "modification_date" : "2022-06-14T15:54:28Z",
       "name" : "needinfo",
       "requestee" : "user3@foobarcorp.example.com",
       "setter" : "user3@foobarcorp.example.com",
       "status" : "?",
       "type_id" : 4
    },
    {
       "creation_date" : "2019-03-27T13:50:29Z",
       "id" : 201663,
       "modification_date" : "2019-03-27T13:50:29Z",
       "name" : "SHIP_STOPPER",
       "requestee" : "user1@foobarcorp.example.com",
       "setter" : "user1@foobarcorp.example.com",
       "status" : "?",
       "type_id" : 2
    },
    {
       "creation_date" : "2022-04-26T08:05:37Z",
       "id" : 264345,
       "modification_date" : "2022-04-26T08:05:37Z",
       "name" : "CCB_Review",
       "setter" : "user2@foobarcorp.example.com",
       "status" : "+",
       "type_id" : 3
    }
 ],
 "groups" : [
    "foobarcorponly",
    "FOOBARCorp Enterprise Partner"
 ],
 "id" : 1047068,
 "is_cc_accessible" : true,
 "is_confirmed" : true,
 "is_creator_accessible" : true,
 "is_open" : true,
 "keywords" : [
    "TRETA",
    "TRETA_ADDRESSED"
 ],
 "last_change_time" : "2023-04-12T01:02:03Z",
 "op_sys" : "Other",
 "platform" : "Other",
 "priority" : "P2 - High",
 "product" : "Enterprise Frobnicator 9000.1",
 "qa_contact" : "user1@foobarcorp.example.com",
 "qa_contact_detail" : {
    "email" : "user1@foobarcorp.example.com",
    "id" : 63803,
    "name" : "user1@foobarcorp.example.com",
    "real_name" : "Firstname1 LastName1"
 },
 "remaining_time" : 0,
 "resolution" : "",
 "see_also" : [],
 "severity" : "Major",
 "status" : "REOPENED",
 "summary" : "L4: test cloud bug123",
 "target_milestone" : "FROB90001Maint-Upd",
 "update_token" : "1683306765-PMQ3v1SB5rHQwTPnDeSPrCAmChAk5itzZn7A_WfGgq4",
 "url" : "https://xxxxxx.foobarcorp.example.com/incident/9999999",
 "version" : "FROB90001Maint-Upd",
 "whiteboard" : "openTreta:1234",
 "attachments": [],
 "comments" : [
    {
       "attachment_id" : null,
       "bug_id" : 1047068,
       "count" : 0,
       "creation_time" : "2017-07-03T13:29:15Z",
       "creator" : "user1@foobarcorp.example.com",
       "id" : 7315202,
       "is_private" : false,
       "tags" : [],
       "text" : "Some comment",
       "time" : "2017-07-03T13:29:15Z"
    },
    {
       "attachment_id" : null,
       "bug_id" : 1047068,
       "count" : 1,
       "creation_time" : "2017-07-03T13:31:23Z",
       "creator" : "bot1@foobarcorp.example.com",
       "id" : 7315205,
       "is_private" : true,
       "tags" : [],
       "text" : "Second comment",
       "time" : "2017-07-03T13:31:23Z"
    }
 ]
}
`

func (cs *clientSuite) TestGetFromJSON(c *C) {
	bz := makeClient("http://bz.foobarcorp.example.com")
	bug, err := bz.GetBugFromJSON(strings.NewReader(sampleJSON))
	c.Assert(err, IsNil)
	c.Assert(bug, NotNil)
	c.Check(bug.Summary, Equals, "L4: test cloud bug123")
	c.Check(bug.AssignedToDetail.RealName, Equals, "Firstname1 LastName1")
	c.Check(bug.Severity, Equals, "Major")
	c.Check(bug.Whiteboard, Equals, "openTreta:1234")
	c.Check(len(bug.Comments), Equals, 2)
	c.Check(bug.Comments[0].Text, Equals, "Some comment")
	c.Check(bug.Comments[0].CreationTime, Equals, time.Date(2017, 07, 03, 13, 29, 15, 0, time.UTC))
	c.Check(bug.CreationTime, Equals, time.Date(2017, 07, 03, 13, 29, 15, 0, time.UTC))
	c.Check(bug.CC, DeepEquals, []string{
		"561726581864@foobarcorp.example.com",
		"user2@foobarcorp.example.com",
		"user1@foobarcorp.example.com",
		"user3@foobarcorp.example.com",
	})
}

func (cs *clientSuite) TestUpdateConnectionClosed(c *C) {
	var ts0 *httptest.Server
	ts0 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts0.CloseClientConnections()
	}))
	defer ts0.Close()
	bz := makeClient(ts0.URL)
	changes := bugzilla.Changes{AddComment: "Some comment"}
	_, err := bz.Update(101234, changes)
	c.Assert(err, ErrorMatches, ".*EOF.*")
}

func (cs *clientSuite) TestDownloadAttachment(c *C) {
	ts0 := cs.makeBugzillaRestServer(c, 1047068)
	defer ts0.Close()
	bz := makeClient(ts0.URL)

	attd, reader, err := bz.DownloadAttachment(766288)
	c.Assert(err, IsNil)
	c.Assert(attd, NotNil)
	c.Assert(reader, NotNil)
	data, err := ioutil.ReadAll(reader)
	c.Assert(data, NotNil)
	c.Assert(err, IsNil)
	decoded, err := attd.DataFromDownload(data)
	c.Assert(err, IsNil)
	c.Assert(string(decoded), Equals, "a\n")
}

func (cs *clientSuite) TestGetAttachment(c *C) {
	ts0 := cs.makeBugzillaRestServer(c, 1047068)
	defer ts0.Close()
	bz := makeClient(ts0.URL)

	att, err := bz.GetAttachment(766288)
	c.Assert(err, IsNil)
	c.Assert(att, NotNil)
	c.Assert(string(att.Data), Equals, "a\n")
}

func (cs *clientSuite) TestUploadAttachment(c *C) {
	ts0, queries, _, processBug := makeBugzillaServerWithChannels()
	defer ts0.Close()
	bz := makeClient(ts0.URL)

	processBug <- `{"ids":[866923]}`

	att := &bugzilla.PostAttachment{
		Data:        []byte("a\n"),
		Summary:     "some summary",
		Filename:    "filename.txt",
		ContentType: "text/plain",
	}
	id, err := bz.UploadAttachment(1047068, att)

	c.Assert(err, IsNil)
	c.Assert(id, Equals, 866923)

	query := <-queries
	ja := jsonassert.New(c)
	ja.Assertf(query, `{"ids":[1047068],"data":"YQo=","file_name": "filename.txt", "content_type": "text/plain","summary": "some summary"}`)
}
