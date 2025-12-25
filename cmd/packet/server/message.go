package server

import "net/http"

// serveGetNewMessage handles GET /new-message requests, which are sent from
// the Create New Message dialog in the GUI.  They have dir= and tag= params,
// with tag= specifying the create tag of the message type to create.  This
// semantically should be a POST, because it's an action, but window.open()
// issues a GET, so that's what it is.  This function creates a new draft
// message of the correct type in the incident, saves it, and then redirects to
// an edit URL to edit it.  That way refreshes of the URL don't create multiple
// messages.
func (s *Server) serveGetNewMessage(w http.ResponseWriter, r *http.Request) {}
