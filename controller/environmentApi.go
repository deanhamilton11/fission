/*
Copyright 2016 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	k_api "k8s.io/client-go/1.5/pkg/api"

	"github.com/fission/fission"
	"github.com/fission/fission/tpr"
)

func (api *API) EnvironmentApiList(w http.ResponseWriter, r *http.Request) {
	envList, err := api.environment.List(k_api.ListOptions{})
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	envs, err := tpr.EnvironmentListFromTPR(envList)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	resp, err := json.Marshal(envs)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	api.respondWithSuccess(w, resp)
}

func (api *API) EnvironmentApiCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	var env fission.Environment
	err = json.Unmarshal(body, &env)
	if err != nil {
		log.Printf("Failed to unmarshal request body: [%v]", body)
		api.respondWithError(w, err)
		return
	}

	t, err := api.environment.Create(tpr.EnvironmentToTPR(&env, nil))
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	m, err := tpr.MetadataFromTPR(&t.Metadata)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	resp, err := json.Marshal(m)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	api.respondWithSuccess(w, resp)
}

func (api *API) EnvironmentApiGet(w http.ResponseWriter, r *http.Request) {
	var m fission.Metadata

	vars := mux.Vars(r)
	m.Name = vars["environment"]
	m.Uid = r.FormValue("uid")         // empty if absent
	m.Version = r.FormValue("version") // empty if absent

	env, err := api.environment.Get(xxx)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	resp, err := json.Marshal(env)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	api.respondWithSuccess(w, resp)
}

func (api *API) EnvironmentApiUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["environment"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	var env fission.Environment
	err = json.Unmarshal(body, &env)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	if name != env.Metadata.Name {
		err = fission.MakeError(fission.ErrorInvalidArgument, "Environment name doesn't match URL")
		api.respondWithError(w, err)
		return
	}

	uid, err := api.EnvironmentStore.Update(&env)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	m := &fission.Metadata{Name: env.Metadata.Name, Uid: uid}
	resp, err := json.Marshal(m)
	if err != nil {
		api.respondWithError(w, err)
		return
	}
	api.respondWithSuccess(w, resp)
}

func (api *API) EnvironmentApiDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var m fission.Metadata
	m.Name = vars["environment"]

	m.Uid = r.FormValue("uid") // empty if uid is absent
	if len(m.Uid) == 0 {
		log.WithFields(log.Fields{"httpTrigger": m.Name}).Info("Deleting all versions")
	}

	err := api.EnvironmentStore.Delete(m)
	if err != nil {
		api.respondWithError(w, err)
		return
	}

	api.respondWithSuccess(w, []byte(""))
}
