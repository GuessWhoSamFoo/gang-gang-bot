package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleGroup_ToggleRole(t *testing.T) {
	cases := []struct {
		name      string
		roleGroup *RoleGroup
		fieldName FieldType
		user      string
		expected  *RoleGroup
	}{
		{
			name:      "should accept",
			roleGroup: NewDefaultRoleGroup(),
			fieldName: AcceptedField,
			user:      "hello",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"hello"},
						Count:     1,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{},
					},
					{
						Icon:      TentativeIcon,
						FieldName: TentativeField,
						Users:     []string{},
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should accept and remove tentative",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{},
					},
					{
						Icon:      TentativeIcon,
						FieldName: TentativeField,
						Users:     []string{"hello"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "hello",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"hello"},
						Count:     1,
					},
					{
						Icon:      TentativeIcon,
						FieldName: TentativeField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should remove accept",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"hello"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "hello",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should accept with multiple",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "baz",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar", "baz"},
						Count:     3,
					},
				},
			},
		},
		{
			name: "should accept with waitlist",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
			fieldName: AcceptedField,
			user:      "baz",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"baz"},
						Count:     1,
					},
				},
			},
		},
		{
			name: "should decline",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo"},
						Count:     1,
						Limit:     1,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{},
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
			fieldName: DeclinedField,
			user:      "bar",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo"},
						Count:     1,
						Limit:     1,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{"bar"},
						Count:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should bump next in waitlist",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"baz"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "baz",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should remove from waitlist",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo"},
						Count:     1,
						Limit:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"bar"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "bar",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo"},
						Count:     1,
						Limit:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should remove from waitlist and move to decline",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{},
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"baz"},
						Count:     1,
					},
				},
			},
			fieldName: DeclinedField,
			user:      "baz",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{"baz"},
						Count:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should bump user in waitlist",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo"},
						Count:     1,
						Limit:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"bar"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "foo",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"bar"},
						Count:     1,
						Limit:     1,
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
		{
			name: "should bump user in waitlist multiple",
			roleGroup: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "bar"},
						Count:     2,
						Limit:     2,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{},
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{"baz"},
						Count:     1,
					},
				},
			},
			fieldName: AcceptedField,
			user:      "bar",
			expected: &RoleGroup{
				Roles: []*Role{
					{
						Icon:      AcceptedIcon,
						FieldName: AcceptedField,
						Users:     []string{"foo", "baz"},
						Count:     2,
						Limit:     2,
					},
					{
						Icon:      DeclinedIcon,
						FieldName: DeclinedField,
						Users:     []string{},
					},
				},
				Waitlist: map[FieldType]*Role{
					AcceptedField: {
						Icon:      "",
						FieldName: WaitlistField,
						Users:     []string{},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.roleGroup.ToggleRole(tc.fieldName, tc.user)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tc.roleGroup)
		})
	}
}