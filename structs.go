package letters

import (
	"fmt"
	"net/mail"
	"time"
)

type void struct{}

var member void

// A set of headers supported directly in letters.structs.Email.Headers
// (and not in letters.structs.Email.Headers.ExtraHeaders)
var knownHeaders = map[string]void{
	"Date":                      member,
	"Sender":                    member,
	"From":                      member,
	"Reply-To":                  member,
	"To":                        member,
	"Cc":                        member,
	"Bcc":                       member,
	"Message-Id":                member,
	"In-Reply-To":               member,
	"References":                member,
	"Subject":                   member,
	"Comments":                  member,
	"Keywords":                  member,
	"Resent-Date":               member,
	"Resent-From":               member,
	"Resent-Sender":             member,
	"Resent-To":                 member,
	"Resent-Cc":                 member,
	"Resent-Bcc":                member,
	"Resent-Message-Id":         member,
	"Content-Transfer-Encoding": member,
	"Content-Type":              member,
	"Content-Disposition":       member,
}

type ContentDisposition string

const (
	ContentDispositionAttachment ContentDisposition = "attachment"
	ContentDispositionInline     ContentDisposition = "inline"
)

var cdMap = map[string]ContentDisposition{
	"attachment": ContentDispositionAttachment,
	"inline":     ContentDispositionInline,
}

const contentTypeMultipartPrefix = "multipart/"

// const contentTypeMultipartAlternative = "multipart/alternative"
// const contentTypeMultipartDigest = "multipart/digest"
const (
	contentTypeMultipartMixed    = "multipart/mixed"
	contentTypeMultipartParallel = "multipart/parallel"
	contentTypeMultipartRelated  = "multipart/related"
)

// const contentTypeMultipartReport = "multipart/report"

// const contentTypeMultipartSigned = "multipart/signed"
// const contentTypeMultipartEncrypted = "multipart/encrypted"

const (
	contentTypeTextPlain    = "text/plain"
	contentTypeTextEnriched = "text/enriched"
	contentTypeTextHtml     = "text/html"
)

type ContentTransferEncoding string

const (
	cte7bit            ContentTransferEncoding = "7bit"
	cte8bit            ContentTransferEncoding = "8bit"
	cteBinary          ContentTransferEncoding = "binary"
	cteQuotedPrintable ContentTransferEncoding = "quoted-printable"
	cteBase64          ContentTransferEncoding = "base64"
)

var cteMap = map[string]ContentTransferEncoding{
	"7bit":             cte7bit,
	"8bit":             cte8bit,
	"binary":           cteBinary,
	"quoted-printable": cteQuotedPrintable,
	"base64":           cteBase64,
}

type UnknownContentTypeError struct {
	contentType string
}

func (e *UnknownContentTypeError) Error() string {
	return fmt.Sprintf("unknown Content-Type %q", e.contentType)
}

type MessageId string

type ContentTypeHeader struct {
	ContentType string
	Params      map[string]string
}

type ContentDispositionHeader struct {
	ContentDisposition ContentDisposition
	Params             map[string]string
}

type Headers struct {
	// RFC 3522 3.6.1.  The Origination Date Field
	// The origination date field consists of the field name "Date" followed
	// by a date-time specification.
	//
	// orig-date       =   "Date:" date-time CRLF
	//
	// The origination date specifies the date and time at which the creator
	// of the message indicated that the message was complete and ready to
	// enter the mail delivery system.  For instance, this might be the time
	// that a user pushes the "send" or "submit" button in an application
	// program.  In any case, it is specifically not intended to convey the
	// time that the message is actually transported, but rather the time at
	// which the human or other creator of the message has put the message
	// into its final form, ready for transport.  (For example, a portable
	// computer user who is not connected to a network might queue a message
	// for delivery.  The origination date is intended to contain the date
	// and time that the user queued the message, not the time when the user
	// connected to the network to send the message.)
	Date time.Time

	// RFC 3522 3.6.2.  Originator Fields
	//
	// The originator fields of a message consist of the from field, the
	// sender field (when applicable), and optionally the reply-to field.
	// The from field consists of the field name "From" and a comma-
	// separated list of one or more mailbox specifications.  If the from
	// field contains more than one mailbox specification in the mailbox-
	// list, then the sender field, containing the field name "Sender" and a
	// single mailbox specification, MUST appear in the message.  In either
	// case, an optional reply-to field MAY also be included, which contains
	// the field name "Reply-To" and a comma-separated list of one or more
	// addresses.
	//
	// from            =   "From:" mailbox-list CRLF
	//
	// sender          =   "Sender:" mailbox CRLF
	//
	// reply-to        =   "Reply-To:" address-list CRLF
	//
	// The originator fields indicate the mailbox(es) of the source of the
	// message.  The "From:" field specifies the author(s) of the message,
	// that is, the mailbox(es) of the person(s) or system(s) responsible
	// for the writing of the message.  The "Sender:" field specifies the
	// mailbox of the agent responsible for the actual transmission of the
	// message.  For example, if a secretary were to send a message for
	// another person, the mailbox of the secretary would appear in the
	// "Sender:" field and the mailbox of the actual author would appear in
	// the "From:" field.  If the originator of the message can be indicated
	// by a single mailbox and the author and transmitter are identical, the
	// "Sender:" field SHOULD NOT be used.  Otherwise, both fields SHOULD
	// appear.
	//
	//    Note: The transmitter information is always present.  The absence
	//    of the "Sender:" field is sometimes mistakenly taken to mean that
	//    the agent responsible for transmission of the message has not been
	//    specified.  This absence merely means that the transmitter is
	//    identical to the author and is therefore not redundantly placed
	//    into the "Sender:" field.
	//
	// The originator fields also provide the information required when
	// replying to a message.  When the "Reply-To:" field is present, it
	// indicates the address(es) to which the author of the message suggests
	// that replies be sent.  In the absence of the "Reply-To:" field,
	// replies SHOULD by default be sent to the mailbox(es) specified in the
	// "From:" field unless otherwise specified by the person composing the
	// reply.
	//
	// In all cases, the "From:" field SHOULD NOT contain any mailbox that
	// does not belong to the author(s) of the message.  See also section
	// 3.6.3 for more information on forming the destination addresses for a
	// reply.
	Sender  *mail.Address
	From    []*mail.Address
	ReplyTo []*mail.Address

	// RFC 3522 3.6.3.  Destination Address Fields
	//
	// The destination fields of a message consist of three possible fields,
	// each of the same form: the field name, which is either "To", "Cc", or
	// "Bcc", followed by a comma-separated list of one or more addresses
	// (either mailbox or group syntax).
	//
	// to              =   "To:" address-list CRLF
	//
	// cc              =   "Cc:" address-list CRLF
	//
	// bcc             =   "Bcc:" [address-list / CFWS] CRLF
	//
	// The destination fields specify the recipients of the message.  Each
	// destination field may have one or more addresses, and the addresses
	// indicate the intended recipients of the message.  The only difference
	// between the three fields is how each is used.
	//
	// The "To:" field contains the address(es) of the primary recipient(s)
	// of the message.
	To  []*mail.Address
	Cc  []*mail.Address
	Bcc []*mail.Address

	// RFC 3522 3.6.4.  Identification Fields
	//
	// Though listed as optional in the table in section 3.6, every message
	// SHOULD have a "Message-ID:" field.  Furthermore, reply messages
	// SHOULD have "In-Reply-To:" and "References:" fields as appropriate
	// and as described below.
	//
	// The "Message-ID:" field contains a single unique message identifier.
	// The "References:" and "In-Reply-To:" fields each contain one or more
	// unique message identifiers, optionally separated by CFWS.
	//
	// The message identifier (msg-id) syntax is a limited version of the
	// addr-spec construct enclosed in the angle bracket characters, "<" and
	// ">".  Unlike addr-spec, this syntax only permits the dot-atom-text
	// form on the left-hand side of the "@" and does not have internal CFWS
	// anywhere in the message identifier.
	//
	//    Note: As with addr-spec, a liberal syntax is given for the right-
	//    hand side of the "@" in a msg-id.  However, later in this section,
	//    the use of a domain for the right-hand side of the "@" is
	//    RECOMMENDED.  Again, the syntax of domain constructs is specified
	//    by and used in other protocols (e.g., [RFC1034], [RFC1035],
	//    [RFC1123], [RFC5321]).  It is therefore incumbent upon
	//    implementations to conform to the syntax of addresses for the
	//    context in which they are used.
	//
	// message-id      =   "Message-ID:" msg-id CRLF
	//
	// in-reply-to     =   "In-Reply-To:" 1*msg-id CRLF
	//
	// references      =   "References:" 1*msg-id CRLF
	//
	// msg-id          =   [CFWS] "<" id-left "@" id-right ">" [CFWS]
	//
	// id-left         =   dot-atom-text / obs-id-left
	//
	// id-right        =   dot-atom-text / no-fold-literal / obs-id-right
	//
	// no-fold-literal =   "[" *dtext "]"
	//
	// The "Message-ID:" field provides a unique message identifier that
	// refers to a particular version of a particular message.  The
	// uniqueness of the message identifier is guaranteed by the host that
	// generates it (see below).  This message identifier is intended to be
	// machine readable and not necessarily meaningful to humans.  A message
	// identifier pertains to exactly one version of a particular message;
	// subsequent revisions to the message each receive new message
	// identifiers.
	//
	//    Note: There are many instances when messages are "changed", but
	//    those changes do not constitute a new instantiation of that
	//    message, and therefore the message would not get a new message
	//    identifier.  For example, when messages are introduced into the
	//    transport system, they are often prepended with additional header
	//    fields such as trace fields (described in section 3.6.7) and
	//    resent fields (described in section 3.6.6).  The addition of such
	//    header fields does not change the identity of the message and
	//    therefore the original "Message-ID:" field is retained.  In all
	//    cases, it is the meaning that the sender of the message wishes to
	//    convey (i.e., whether this is the same message or a different
	//    message) that determines whether or not the "Message-ID:" field
	//    changes, not any particular syntactic difference that appears (or
	//    does not appear) in the message.
	//
	// The "In-Reply-To:" and "References:" fields are used when creating a
	// reply to a message.  They hold the message identifier of the original
	// message and the message identifiers of other messages (for example,
	// in the case of a reply to a message that was itself a reply).  The
	// "In-Reply-To:" field may be used to identify the message (or
	// messages) to which the new message is a reply, while the
	// "References:" field may be used to identify a "thread" of
	// conversation.
	//
	// When creating a reply to a message, the "In-Reply-To:" and
	// "References:" fields of the resultant message are constructed as
	// follows:
	//
	// The "In-Reply-To:" field will contain the contents of the
	// "Message-ID:" field of the message to which this one is a reply (the
	// "parent message").  If there is more than one parent message, then
	// the "In-Reply-To:" field will contain the contents of all of the
	// parents' "Message-ID:" fields.  If there is no "Message-ID:" field in
	// any of the parent messages, then the new message will have no "In-
	// Reply-To:" field.
	//
	// The "References:" field will contain the contents of the parent's
	// "References:" field (if any) followed by the contents of the parent's
	// "Message-ID:" field (if any).  If the parent message does not contain
	// a "References:" field but does have an "In-Reply-To:" field
	// containing a single message identifier, then the "References:" field
	// will contain the contents of the parent's "In-Reply-To:" field
	// followed by the contents of the parent's "Message-ID:" field (if
	// any).  If the parent has none of the "References:", "In-Reply-To:",
	// or "Message-ID:" fields, then the new message will have no
	// "References:" field.
	//
	//    Note: Some implementations parse the "References:" field to
	//    display the "thread of the discussion".  These implementations
	//    assume that each new message is a reply to a single parent and
	//    hence that they can walk backwards through the "References:" field
	//    to find the parent of each message listed there.  Therefore,
	//    trying to form a "References:" field for a reply that has multiple
	//    parents is discouraged; how to do so is not defined in this
	//    document.
	//
	// The message identifier (msg-id) itself MUST be a globally unique
	// identifier for a message.  The generator of the message identifier
	// MUST guarantee that the msg-id is unique.  There are several
	// algorithms that can be used to accomplish this.  Since the msg-id has
	// a similar syntax to addr-spec (identical except that quoted strings,
	// comments, and folding white space are not allowed), a good method is
	// to put the domain name (or a domain literal IP address) of the host
	// on which the message identifier was created on the right-hand side of
	// the "@" (since domain names and IP addresses are normally unique),
	// and put a combination of the current absolute date and time along
	// with some other currently unique (perhaps sequential) identifier
	// available on the system (for example, a process id number) on the
	// left-hand side.  Though other algorithms will work, it is RECOMMENDED
	// that the right-hand side contain some domain identifier (either of
	// the host itself or otherwise) such that the generator of the message
	// identifier can guarantee the uniqueness of the left-hand side within
	// the scope of that domain.
	//
	// Semantically, the angle bracket characters are not part of the
	// msg-id; the msg-id is what is contained between the two angle bracket
	// characters.
	MessageID  MessageId
	InReplyTo  []MessageId
	References []MessageId

	// RFC 3522 3.6.5.  Informational Fields
	//
	// The informational fields are all optional.  The "Subject:" and
	// "Comments:" fields are unstructured fields as defined in section
	// 2.2.1, and therefore may contain text or folding white space.  The
	// "Keywords:" field contains a comma-separated list of one or more
	// words or quoted-strings.
	//
	// subject         =   "Subject:" unstructured CRLF
	//
	// comments        =   "Comments:" unstructured CRLF
	//
	// keywords        =   "Keywords:" phrase *("," phrase) CRLF
	//
	// These three fields are intended to have only human-readable content
	// with information about the message.  The "Subject:" field is the most
	// common and contains a short string identifying the topic of the
	// message.  When used in a reply, the field body MAY start with the
	// string "Re: " (an abbreviation of the Latin "in re", meaning "in the
	// matter of") followed by the contents of the "Subject:" field body of
	// the original message.  If this is done, only one instance of the
	// literal string "Re: " ought to be used since use of other strings or
	// more than one instance can lead to undesirable consequences.  The
	// "Comments:" field contains any additional comments on the text of the
	// body of the message.  The "Keywords:" field contains a comma-
	// separated list of important words and phrases that might be useful
	// for the recipient.
	Subject  string
	Comments string
	Keywords []string

	// RFC 3522 3.6.6.  Resent Fields
	//
	// Resent fields SHOULD be added to any message that is reintroduced by
	// a user into the transport system.  A separate set of resent fields
	// SHOULD be added each time this is done.  All of the resent fields
	// corresponding to a particular resending of the message SHOULD be
	// grouped together.  Each new set of resent fields is prepended to the
	// message; that is, the most recent set of resent fields appears
	// earlier in the message.  No other fields in the message are changed
	// when resent fields are added.
	//
	// Each of the resent fields corresponds to a particular field elsewhere
	// in the syntax.  For instance, the "Resent-Date:" field corresponds to
	// the "Date:" field and the "Resent-To:" field corresponds to the "To:"
	// field.  In each case, the syntax for the field body is identical to
	// the syntax given previously for the corresponding field.
	//
	// When resent fields are used, the "Resent-From:" and "Resent-Date:"
	// fields MUST be sent.  The "Resent-Message-ID:" field SHOULD be sent.
	// "Resent-Sender:" SHOULD NOT be used if "Resent-Sender:" would be
	// identical to "Resent-From:".
	//
	// resent-date     =   "Resent-Date:" date-time CRLF
	//
	// resent-from     =   "Resent-From:" mailbox-list CRLF
	//
	// resent-sender   =   "Resent-Sender:" mailbox CRLF
	//
	// resent-to       =   "Resent-To:" address-list CRLF
	//
	// resent-cc       =   "Resent-Cc:" address-list CRLF
	//
	// resent-bcc      =   "Resent-Bcc:" [address-list / CFWS] CRLF
	//
	// resent-msg-id   =   "Resent-Message-ID:" msg-id CRLF
	//
	// Resent fields are used to identify a message as having been
	// reintroduced into the transport system by a user.  The purpose of
	// using resent fields is to have the message appear to the final
	// recipient as if it were sent directly by the original sender, with
	// all of the original fields remaining the same.  Each set of resent
	// fields correspond to a particular resending event.  That is, if a
	// message is resent multiple times, each set of resent fields gives
	// identifying information for each individual time.  Resent fields are
	// strictly informational.  They MUST NOT be used in the normal
	// processing of replies or other such automatic actions on messages.
	//
	//    Note: Reintroducing a message into the transport system and using
	//    resent fields is a different operation from "forwarding".
	//    "Forwarding" has two meanings: One sense of forwarding is that a
	//    mail reading program can be told by a user to forward a copy of a
	//    message to another person, making the forwarded message the body
	//    of the new message.  A forwarded message in this sense does not
	//    appear to have come from the original sender, but is an entirely
	//    new message from the forwarder of the message.  Forwarding may
	//    also mean that a mail transport program gets a message and
	//    forwards it on to a different destination for final delivery.
	//    Resent header fields are not intended for use with either type of
	//    forwarding.
	//
	// The resent originator fields indicate the mailbox of the person(s) or
	// system(s) that resent the message.  As with the regular originator
	// fields, there are two forms: a simple "Resent-From:" form, which
	// contains the mailbox of the individual doing the resending, and the
	// more complex form, when one individual (identified in the "Resent-
	// Sender:" field) resends a message on behalf of one or more others
	// (identified in the "Resent-From:" field).
	//
	//    Note: When replying to a resent message, replies behave just as
	//    they would with any other message, using the original "From:",
	//    "Reply-To:", "Message-ID:", and other fields.  The resent fields
	//    are only informational and MUST NOT be used in the normal
	//    processing of replies.
	//
	// The "Resent-Date:" indicates the date and time at which the resent
	// message is dispatched by the resender of the message.  Like the
	// "Date:" field, it is not the date and time that the message was
	// actually transported.
	//
	// The "Resent-To:", "Resent-Cc:", and "Resent-Bcc:" fields function
	// identically to the "To:", "Cc:", and "Bcc:" fields, respectively,
	// except that they indicate the recipients of the resent message, not
	// the recipients of the original message.
	//
	// The "Resent-Message-ID:" field provides a unique identifier for the
	// resent message.
	ResentDate      time.Time
	ResentFrom      []*mail.Address
	ResentSender    *mail.Address
	ResentTo        []*mail.Address
	ResentCc        []*mail.Address
	ResentBcc       []*mail.Address
	ResentMessageID MessageId

	// RFC 2045 5.  Content-Type Header Field
	//
	// The purpose of the Content-Type field is to describe the data
	// contained in the body fully enough that the receiving user agent can
	// pick an appropriate agent or mechanism to present the data to the
	// user, or otherwise deal with the data in an appropriate manner. The
	// value in this field is called a media type.
	//
	// HISTORICAL NOTE:  The Content-Type header field was first defined in
	// RFC 1049.  RFC 1049 used a simpler and less powerful syntax, but one
	// that is largely compatible with the mechanism given here.
	//
	// The Content-Type header field specifies the nature of the data in the
	// body of an entity by giving media type and subtype identifiers, and
	// by providing auxiliary information that may be required for certain
	// media types.  After the media type and subtype names, the remainder
	// of the header field is simply a set of parameters, specified in an
	// attribute=value notation.  The ordering of parameters is not
	// significant.
	//
	// In general, the top-level media type is used to declare the general
	// type of data, while the subtype specifies a specific format for that
	// type of data.  Thus, a media type of "image/xyz" is enough to tell a
	// user agent that the data is an image, even if the user agent has no
	// knowledge of the specific image format "xyz".  Such information can
	// be used, for example, to decide whether or not to show a user the raw
	// data from an unrecognized subtype -- such an action might be
	// reasonable for unrecognized subtypes of text, but not for
	// unrecognized subtypes of image or audio.  For this reason, registered
	// subtypes of text, image, audio, and video should not contain embedded
	// information that is really of a different type.  Such compound
	// formats should be represented using the "multipart" or "application"
	// types.
	//
	// Parameters are modifiers of the media subtype, and as such do not
	// fundamentally affect the nature of the content.  The set of
	// meaningful parameters depends on the media type and subtype.  Most
	// parameters are associated with a single specific subtype.  However, a
	// given top-level media type may define parameters which are applicable
	// to any subtype of that type.  Parameters may be required by their
	// defining content type or subtype or they may be optional. MIME
	// implementations must ignore any parameters whose names they do not
	// recognize.
	//
	// For example, the "charset" parameter is applicable to any subtype of
	// "text", while the "boundary" parameter is required for any subtype of
	// the "multipart" media type.
	//
	// There are NO globally-meaningful parameters that apply to all media
	// types.  Truly global mechanisms are best addressed, in the MIME
	// model, by the definition of additional Content-* header fields.
	//
	// An initial set of seven top-level media types is defined in RFC 2046.
	// Five of these are discrete types whose content is essentially opaque
	// as far as MIME processing is concerned.  The remaining two are
	// composite types whose contents require additional handling by MIME
	// processors.
	//
	// This set of top-level media types is intended to be substantially
	// complete.  It is expected that additions to the larger set of
	// supported types can generally be accomplished by the creation of new
	// subtypes of these initial types.  In the future, more top-level types
	// may be defined only by a standards-track extension to this standard.
	// If another top-level type is to be used for any reason, it must be
	// given a name starting with "X-" to indicate its non-standard status
	// and to avoid a potential conflict with a future official name.
	ContentType        ContentTypeHeader
	ContentDisposition ContentDispositionHeader
	ExtraHeaders       map[string][]string
}

type emailBodies struct {
	text string

	enrichedText string // See RFC 1523, RFC 1563, and RFC 1896
	html         string

	InlineFiles   []InlineFile
	AttachedFiles []AttachedFile
}

func (eb *emailBodies) extend(b emailBodies) {
	eb.text += b.text
	eb.enrichedText += b.enrichedText
	eb.html += b.html
	eb.InlineFiles = append(eb.InlineFiles, b.InlineFiles...)
	eb.AttachedFiles = append(eb.AttachedFiles, b.AttachedFiles...)
}

type Email struct {
	Headers Headers

	Text         string
	EnrichedText string // See RFC 1523, RFC 1563, and RFC 1896
	HTML         string

	InlineFiles   []InlineFile
	AttachedFiles []AttachedFile
}

type InlineFile struct {
	ContentID          string
	ContentType        ContentTypeHeader
	ContentDisposition ContentDispositionHeader
	Data               []byte
}

type AttachedFile struct {
	ContentType        ContentTypeHeader
	ContentDisposition ContentDispositionHeader
	Data               []byte
}
