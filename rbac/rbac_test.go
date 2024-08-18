package rbac

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/TechMaster/core/pmodel"
	"github.com/stretchr/testify/assert"
)

type Test[T interface{}, E interface{}] struct {
	ID     int
	Name   string
	Want   bool
	Expect E
	Args   T
}

func TestLoadRoles(t *testing.T) {
	dataTest := []Test[[]pmodel.Role, any]{
		{
			ID: 1, Name: "Case 1", Want: true,
			Args: []pmodel.Role{
				{ID: 1, Name: "admin"},
				{ID: 2, Name: "student"},
				{ID: 3, Name: "trainer"},
				{ID: 4, Name: "sale"},
				{ID: 5, Name: "employer"},
				{ID: 6, Name: "author"},
				{ID: 7, Name: "editor"},
				{ID: 8, Name: "maintainer"},
			},
			Expect: map[string]int{"admin": 1, "student": 2, "trainer": 3, "sale": 4, "employer": 5, "author": 6, "editor": 7, "maintainer": 8},
		},
		{
			ID: 2, Name: "Case 2", Want: false,
			Args: []pmodel.Role{
				{ID: 1, Name: "admin"},
				{ID: 2, Name: "student"},
				{ID: 3, Name: "trainer"},
				{ID: 4, Name: "sale"},
				{ID: 5, Name: "employer"},
				{ID: 6, Name: "author"},
				{ID: 7, Name: "editor"},
				{ID: 8, Name: "maintainer"},
			},
			Expect: map[string]int{"admin": 1, "student": 2, "trainer": 3, "sale": 4},
		},
		{
			ID: 1, Name: "Case 3", Want: false,
			Args: []pmodel.Role{
				{ID: 1, Name: "admin"},
				{ID: 2, Name: "trainer"},
				{ID: 3, Name: "student"},
				{ID: 4, Name: "sale"},
				{ID: 5, Name: "employer"},
				{ID: 6, Name: "author"},
				{ID: 7, Name: "editor"},
				{ID: 8, Name: "maintainer"},
			},
			Expect: map[string]int{"admin": 1, "student": 2, "trainer": 3, "sale": 4, "employer": 5, "author": 6, "editor": 7, "maintainer": 8},
		},
	}
	for _, item := range dataTest {
		t.Run(item.Name, func(t *testing.T) {
			LoadRoles(func() []pmodel.Role {
				return item.Args
			})
			result := true
			if len(Roles) != len(item.Expect.(map[string]int)) {
				result = false
			}
			for key, value := range item.Expect.(map[string]int) {
				if Roles[key] != value {
					result = false
					break
				}
			}
			if result != item.Want {
				t.Errorf("Got %v, want %v", result, item.Want)
			}
		})
	}
}

func TestLoadRules(t *testing.T) {
	dataTest := []Test[[]pmodel.Rule, any]{
		{
			ID: 1, Name: "Case 1", Want: true,
			Args: []pmodel.Rule{
				{
					ID:         1,
					Name:       "API",
					Roles:      []int{1, 2, 3},
					AccessType: "allow",
					Method:     "POST",
					Path:       "/api",
					IsPrivate:  true,
				},
			},
			Expect: map[string]Route{
				"POST /api": {
					AccessType: "allow",
					Path:       "/api",
					Method:     "POST",
					IsPrivate:  true,
					Roles: map[int]interface{}{
						1: true,
						2: true,
						3: true,
					},
				},
			},
		},
		{
			ID: 2, Name: "Case 2", Want: false,
			Args: []pmodel.Rule{
				{
					ID:         1,
					Name:       "API",
					Roles:      []int{1, 2, 3},
					AccessType: "allow",
					Method:     "POST",
					Path:       "/api",
					IsPrivate:  true,
				},
			},
			Expect: map[string]Route{
				"POST /api": {
					AccessType: "allow",
					Path:       "/api",
					Method:     "POST",
					IsPrivate:  false,
					Roles:      map[int]interface{}{},
				},
			},
		},
	}

	for _, item := range dataTest {
		t.Run(item.Name, func(t *testing.T) {
			LoadRules(func() []pmodel.Rule {
				return item.Args
			})
			result := true
			for key, value := range item.Expect.(map[string]Route) {
				if routesRoles[key].AccessType != value.AccessType {
					result = false
				}
				if routesRoles[key].Path != value.Path {
					result = false
				}
				if routesRoles[key].Method != value.Method {
					result = false
				}
				if routesRoles[key].IsPrivate != value.IsPrivate {
					result = false
				}
				if reflect.DeepEqual(routesRoles[key].Roles, value.Roles) != true {
					result = false
				}
			}
			if result != item.Want {
				t.Errorf("Got %v, want %v", result, item.Want)
			}
		})
	}
}

func TestCheckUserRouteRoleIntersect(t *testing.T) {
	dataTest := []Test[map[string]map[int]interface{}, bool]{
		{
			ID: 1, Name: "Case 1", Want: true,
			Args: map[string]map[int]interface{}{
				"1": {1: true, 2: true, 3: true},
				"2": {1: true},
			},
			Expect: true,
		},
		{
			ID: 2, Name: "Case 2", Want: true,
			Args: map[string]map[int]interface{}{
				"1": {1: true, 2: true, 3: true},
				"2": {4: true, 5: true},
			},
			Expect: false,
		},
		{
			ID: 3, Name: "Case 3", Want: false,
			Args: map[string]map[int]interface{}{
				"1": {1: true, 2: true, 3: true},
				"2": {1: true, 5: true},
			},
			Expect: false,
		},
	}

	for _, item := range dataTest {
		t.Run(item.Name, func(t *testing.T) {
			result := checkUser_RouteRole_Intersect(item.Args["1"], item.Args["2"])
			r := true
			if result != item.Expect {
				r = false
			}
			if r != item.Want {
				t.Errorf("Got %v, want %v", r, item.Want)
			}
		})
	}
}

func TestCheckAdmin(t *testing.T) {
	dataTest := []Test[map[int]interface{}, bool]{
		{
			ID: 1, Name: "Case 1", Want: true,
			Args: map[int]interface{}{
				1: true,
				2: true,
			},
			Expect: true,
		},
		{
			ID: 2, Name: "Case 2", Want: true,
			Args: map[int]interface{}{
				2: true,
				3: true,
			},
			Expect: false,
		},
		{
			ID: 3, Name: "Case 3", Want: false,
			Args: map[int]interface{}{
				1: true,
				5: true,
			},
			Expect: false,
		},
	}

	for _, item := range dataTest {
		t.Run(item.Name, func(t *testing.T) {
			result := checkAdmin(item.Args)
			r := true
			if result != item.Expect {
				r = false
			}
			if r != item.Want {
				t.Errorf("Got %v, want %v", r, item.Want)
			}
		})
	}
}

func TestAssignRoles(t *testing.T) {
	dataTest := []Test[Route, Route]{
		{
			ID: 1, Name: "Case 1", Want: true,
			Args: Route{
				Path:      "/api",
				Method:    "POST",
				IsPrivate: true,
				Roles: map[int]interface{}{
					1: true,
					2: true,
				},
				AccessType: "allow",
			},
			Expect: Route{
				Path:      "/api",
				Method:    "POST",
				IsPrivate: true,
				Roles: map[int]interface{}{
					1: true,
					2: true,
				},
				AccessType: "allow",
			},
		},
		{
			ID: 1, Name: "Case 2", Want: false,
			Args: Route{
				Path:      "/api",
				Method:    "POST",
				IsPrivate: true,
				Roles: map[int]interface{}{
					1: true,
					2: true,
				},
				AccessType: "allow",
			},
			Expect: Route{
				Path:      "/api",
				Method:    "POST",
				IsPrivate: false,
				Roles: map[int]interface{}{
					1: true,
					2: true,
				},
				AccessType: "allow",
			},
		},
	}

	for _, item := range dataTest {
		t.Run(item.Name, func(t *testing.T) {
			assignRoles(item.Args)
			result := true
			if routesRoles[item.Expect.Method+" "+item.Expect.Path].AccessType != item.Expect.AccessType {
				result = false
			}
			if routesRoles[item.Expect.Method+" "+item.Expect.Path].Path != item.Expect.Path {
				result = false
			}
			if routesRoles[item.Expect.Method+" "+item.Expect.Path].Method != item.Expect.Method {
				result = false
			}
			if routesRoles[item.Expect.Method+" "+item.Expect.Path].IsPrivate != item.Expect.IsPrivate {
				result = false
			}
			if reflect.DeepEqual(routesRoles[item.Expect.Method+" "+item.Expect.Path].Roles, item.Expect.Roles) != true {
				result = false
			}
			if result != item.Want {
				t.Errorf("Got %v, want %v", result, item.Want)
			}
		})
	}
}

func TestAllow(t *testing.T) {
	roles := []int{1, 2, 3}
	expectedRoles := pmodel.Roles{1: true, 2: true, 3: true}
	expectedType := "allow"

	roleExp := Allow(roles...)
	actualRoles, actualType := roleExp()

	assert.Equal(t, expectedRoles, actualRoles, "Roles should match expected roles")
	assert.Equal(t, expectedType, actualType, "Type should be 'allow'")
}

func TestAllowAll(t *testing.T) {
	// Assuming Roles is a map of all possible roles
	expectedRoles := pmodel.Roles{1: true, 2: true, 3: true}
	expectedType := "allow_all"

	roleExp := AllowAll()
	actualRoles, actualType := roleExp()
	for key, value := range expectedRoles {
		t.Run("Case AllowAll "+strconv.Itoa(key), func(t *testing.T) {
			assert.Equal(t, value, actualRoles[key], "Role should match expected role")
			assert.Equal(t, expectedType, actualType, "Type should be 'allow_all'")
		})
	}
}

func TestForbid(t *testing.T) {
	roles := []int{1, 2, 3}
	expectedRoles := pmodel.Roles{1: false, 2: false, 3: false}
	expectedType := "forbid"

	roleExp := Forbid(roles...)
	actualRoles, actualType := roleExp()

	assert.Equal(t, expectedRoles, actualRoles, "Roles should match expected roles")
	assert.Equal(t, expectedType, actualType, "Type should be 'forbid'")
}

func TestForbidAll(t *testing.T) {
	// Assuming Roles is a map of all possible roles
	expectedRoles := pmodel.Roles{1: false, 2: false, 3: false}
	expectedType := "forbid_all"

	roleExp := ForbidAll()
	actualRoles, actualType := roleExp()
	for key, value := range expectedRoles {
		t.Run("Case ForbidAll "+strconv.Itoa(key), func(t *testing.T) {
			assert.Equal(t, value, actualRoles[key], "Role should match expected role")
			assert.Equal(t, expectedType, actualType, "Type should be 'forbid_all'")
		})
	}

}
