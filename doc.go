// Package letters is a minimal Go library for parsing plaintext and MIME
// emails.
//
// It correctly handles text and MIME mime-types, Base64 and Quoted-Printable
// Content-Transfer-Encoding, as well as any text encoding that Go standard
// library is capable of handling. Letters will parse an email into a simple struct
// with standard headers and text, enriched text, and HTML content, and decode
// inline and attached files.
package letters
