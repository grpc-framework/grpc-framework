/*
Copyright 2025 Luis Pabon

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
package role

import (
	"testing"

	authv1 "github.com/grpc-framework/grpc-framework/v2/apis/auth/apiv1"
	"github.com/stretchr/testify/assert"
)

func TestMapRoleStore(t *testing.T) {
	m := NewMapRoleStore()
	assert.NotNil(t, m.store)
}

func TestMapRoleStoreValues(t *testing.T) {
	m := NewMapRoleStoreFromRoles(DefaultRoles)
	assert.NotNil(t, m.store)
	assert.Len(t, m.store, len(DefaultRoles))
	l, err := m.List()
	assert.NoError(t, err)
	assert.Len(t, l, len(DefaultRoles))

	r := DefaultRoles[0]
	v, ok := m.Get(r.Name)
	assert.True(t, ok)
	assert.Equal(t, v, DefaultRoles[0])

	err = m.Set(&authv1.Role{
		Name:  "test",
		Rules: []*authv1.Rule{},
	})
	assert.Len(t, m.store, len(DefaultRoles)+1)
	l, err = m.List()
	assert.NoError(t, err)
	assert.Len(t, l, len(DefaultRoles)+1)

	m.Delete("test")
	assert.Len(t, m.store, len(DefaultRoles))
	l, err = m.List()
	assert.NoError(t, err)
	assert.Len(t, l, len(DefaultRoles))
}
