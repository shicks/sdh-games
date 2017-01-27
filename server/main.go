package server

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"net/http"
)

type Cxn struct {
	Id string
}

func init() {
	http.Handle("/rpc/login", Handler(login))
	http.Handle("/rpc/send", Handler(send))
	http.Handle("/_ah/channel/disconnected/", Handler(disconnect))
}

type loginRequest struct {
	Id string `json:"id"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func login(r *Request) {
	req := loginRequest{}
	r.PostJson(&req)
	if r.Done() {
		return
	}

	tok, err := channel.Create(r.Ctx(), req.Id)
	if err != nil {
		r.Fail("Couldn't create channel", http.StatusInternalServerError)
		r.Ctx().Errorf("channel.Create(%v): %v", req.Id, err)
		return
	}

	err = datastore.RunInTransaction(r.Ctx(), func(c appengine.Context) error {
		k := datastore.NewKey(c, "Cxn", req.Id, 0, nil)
		cxn := new(Cxn)
		cxn.Id = req.Id
		_, err := datastore.Put(c, k, cxn)
		return err
	}, nil)
	if err != nil {
 		r.Fail("Couldn't write id in datastore", http.StatusInternalServerError)
		r.Ctx().Errorf("datastore.Put: %v", err)
		return
	}

	r.RespondJson(loginResponse{tok})
}

type sendRequest struct {
	Text string `json:"text"`
}
type sendResponse struct {}
type sendMessage struct {
	Text string `json:"text"`
}

func send(r *Request) {
	req := sendRequest{}
	r.PostJson(&req)
	if r.Done() {
		return
	}

	q := datastore.NewQuery("Cxn")
	var cxns []Cxn
	_, err := q.GetAll(r.Ctx(), &cxns)
	if err != nil {
		r.Fail("Couldn't query datastore", http.StatusInternalServerError)
		r.Ctx().Errorf("q.GetAll: %v", err)
		return
	}

	msg := sendMessage{req.Text}
	for _, cxn := range cxns {
    r.Ctx().Infof("sending %v <- %v", cxn, msg)
		err := channel.SendJSON(r.Ctx(), cxn.Id, msg)
		if err != nil {
			r.Ctx().Errorf("sending chat: %v", err)
		}
	}

	r.RespondJson(sendResponse{})
}

func disconnect(r *Request) {
	r.Ctx().Infof("DISCONNECT")
	clientId := r.HttpRequest().FormValue("from")
	q := datastore.NewQuery("Cxn").Filter("Id=", clientId).KeysOnly()
	cnt := 0
  for it := q.Run(r.Ctx()); ; {
		k, err := it.Next(nil)
		if err != nil {
			if err != datastore.Done {
				r.Ctx().Errorf("query failed: %v", err)
			}
			break
		}
		err = datastore.Delete(r.Ctx(), k)
		if err != nil {
			r.Ctx().Errorf("failed to delete %v: %v", k, err)
		}
		cnt += 1
	}
	r.Ctx().Infof("Deleted %d channels", cnt)

	// random (dis)connections -> should only clean up connections
	// if they're disconnected for a long time?

}
