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

package tpr

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/watch"
	"k8s.io/client-go/1.5/rest"
)

type (
	HTTPTriggerInterface interface {
		Create(*HTTPTrigger) (*HTTPTrigger, error)
		Get(name string) (*HTTPTrigger, error)
		Update(*HTTPTrigger) (*HTTPTrigger, error)
		Delete(name string, options *api.DeleteOptions) error
		List(opts api.ListOptions) (*HTTPTriggerList, error)
		Watch(opts api.ListOptions) (watch.Interface, error)
	}

	httpTriggerClient struct {
		client    *rest.RESTClient
		namespace string
	}
)

func MakeHTTPTriggerInterface(tprClient *rest.RESTClient, namespace string) HTTPTriggerInterface {
	return &httpTriggerClient{
		client:    tprClient,
		namespace: namespace,
	}
}

func (c *httpTriggerClient) Create(obj *HTTPTrigger) (*HTTPTrigger, error) {
	var result HTTPTrigger
	err := c.client.Post().
		Resource("httptriggers").
		Namespace(c.namespace).
		Body(obj).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *httpTriggerClient) Get(name string) (*HTTPTrigger, error) {
	var result HTTPTrigger
	err := c.client.Get().
		Resource("httptriggers").
		Namespace(c.namespace).
		Name(name).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *httpTriggerClient) Update(obj *HTTPTrigger) (*HTTPTrigger, error) {
	var result HTTPTrigger
	err := c.client.Put().
		Resource("httptriggers").
		Namespace(c.namespace).
		Name(obj.Metadata.Name).
		Body(obj).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *httpTriggerClient) Delete(name string, opts *api.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.namespace).
		Resource("httptriggers").
		Name(name).
		Body(opts).
		Do().
		Error()
}

func (c *httpTriggerClient) List(opts api.ListOptions) (*HTTPTriggerList, error) {
	var result HTTPTriggerList
	err := c.client.Get().
		Namespace(c.namespace).
		Resource("httptriggers").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *httpTriggerClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.namespace).
		Resource("httptriggers").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
