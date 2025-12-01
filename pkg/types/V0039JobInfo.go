// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"

	"k8s.io/utils/ptr"
	"k8s.io/utils/set"

	api "github.com/supergate-hub/slurm-client/api/v0039"
	"github.com/supergate-hub/slurm-client/pkg/object"
	"github.com/supergate-hub/slurm-client/pkg/utils"
)

const (
	ObjectTypeV0039JobInfo = "V0039JobInfo"
)

type V0039JobInfo struct {
	api.V0039JobInfo
}

// GetKey implements Object.
func (o *V0039JobInfo) GetKey() object.ObjectKey {
	return object.ObjectKey(fmt.Sprintf("%d", ptr.Deref(o.JobId, 0)))
}

// GetType implements Object.
func (o *V0039JobInfo) GetType() object.ObjectType {
	return ObjectTypeV0039JobInfo
}

// DeepCopyObject implements Object.
func (o *V0039JobInfo) DeepCopyObject() object.Object {
	return o.DeepCopy()
}

func (o *V0039JobInfo) DeepCopy() *V0039JobInfo {
	out := new(V0039JobInfo)
	utils.RemarshalOrDie(o, out)
	return out
}

func (o *V0039JobInfo) GetStateAsSet() set.Set[string] {
	out := make(set.Set[string])
	state := ptr.Deref(o.JobState, "FAILED")
	out.Insert(state)
	return out
}

type V0039JobInfoList struct {
	Items []V0039JobInfo
}

// GetType implements ObjectList.
func (o *V0039JobInfoList) GetType() object.ObjectType {
	return ObjectTypeV0039JobInfo
}

// GetItems implements ObjectList.
func (o *V0039JobInfoList) GetItems() []object.Object {
	list := make([]object.Object, len(o.Items))
	for i, item := range o.Items {
		list[i] = item.DeepCopyObject()
	}
	return list
}

// AppendItem implements ObjectList.
func (o *V0039JobInfoList) AppendItem(object object.Object) {
	out, ok := object.(*V0039JobInfo)
	if ok {
		utils.RemarshalOrDie(object, out)
		o.Items = append(o.Items, *out)
	}
}

// DeepCopyObjectList implements ObjectList.
func (o *V0039JobInfoList) DeepCopyObjectList() object.ObjectList {
	out := new(V0039JobInfoList)
	out.Items = make([]V0039JobInfo, len(o.Items))
	for i, item := range o.Items {
		out.Items[i] = *item.DeepCopy()
	}
	return out
}
