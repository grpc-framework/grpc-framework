/*
Generirc
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

import authv1 "github.com/grpc-framework/grpc-framework/v2/apis/auth/apiv1"

type MapRoleStore struct {
	store map[string]*authv1.Role
}

func NewMapRoleStore() *MapRoleStore {
	return &MapRoleStore{
		store: make(map[string]*authv1.Role),
	}
}

func NewMapRoleStoreFromRoles(roles []*authv1.Role) *MapRoleStore {
	m := NewMapRoleStore()
	for _, role := range roles {
		m.Set(role)
	}
	return m
}

func (m *MapRoleStore) Set(role *authv1.Role) error {
	m.store[role.Name] = role
	return nil
}

func (m *MapRoleStore) Delete(roleName string) {
	delete(m.store, roleName)
}

func (m *MapRoleStore) Get(roleName string) (*authv1.Role, bool) {
	v, ok := m.store[roleName]
	return v, ok
}

func (m *MapRoleStore) List() ([]*authv1.Role, error) {
	l := make([]*authv1.Role, 0, len(m.store))
	for _, v := range m.store {
		l = append(l, v)
	}
	return l, nil
}
