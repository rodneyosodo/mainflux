// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/things/policies"
)

var (
	_ mainflux.Response = (*identityRes)(nil)
	_ mainflux.Response = (*authorizeRes)(nil)
	_ mainflux.Response = (*addPolicyRes)(nil)
	_ mainflux.Response = (*viewPolicyRes)(nil)
	_ mainflux.Response = (*listPolicyRes)(nil)
	_ mainflux.Response = (*updatePolicyRes)(nil)
	_ mainflux.Response = (*deletePolicyRes)(nil)
)

type pageRes struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
	Total  uint64 `json:"total"`
}

type authorizeRes struct {
	ThingID    string `json:"thing_id"`
	Authorized bool   `json:"authorized"`
}

func (res authorizeRes) Code() int {
	if !res.Authorized {
		return http.StatusForbidden
	}

	return http.StatusOK
}

func (res authorizeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res authorizeRes) Empty() bool {
	return false
}

type identityRes struct {
	ID string `json:"id"`
}

func (res identityRes) Code() int {
	return http.StatusOK
}

func (res identityRes) Headers() map[string]string {
	return map[string]string{}
}

func (res identityRes) Empty() bool {
	return false
}

type addPolicyRes struct {
	created         bool
	policies.Policy `json:",inline"`
}

func (res addPolicyRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusBadRequest
}

func (res addPolicyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res addPolicyRes) Empty() bool {
	return true
}

type viewPolicyRes struct {
	policies.Policy `json:",inline"`
}

func (res viewPolicyRes) Code() int {
	return http.StatusOK
}

func (res viewPolicyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewPolicyRes) Empty() bool {
	return false
}

type updatePolicyRes struct {
	updated         bool
	policies.Policy `json:",inline"`
}

func (res updatePolicyRes) Code() int {
	if res.updated {
		return http.StatusNoContent
	}

	return http.StatusBadRequest
}

func (res updatePolicyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updatePolicyRes) Empty() bool {
	return true
}

type listPolicyRes struct {
	pageRes
	Policies []policies.Policy `json:"policies"`
}

func (res listPolicyRes) Code() int {
	return http.StatusOK
}

func (res listPolicyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res listPolicyRes) Empty() bool {
	return false
}

type deletePolicyRes struct {
	deleted bool
}

func (res deletePolicyRes) Code() int {
	if res.deleted {
		return http.StatusNoContent
	}

	return http.StatusBadRequest
}

func (res deletePolicyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res deletePolicyRes) Empty() bool {
	return true
}
