package route

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

var (
	r1 = &Route{
		&BoxRoute{
			Model: gorm.Model{
				ID: 1,
			},
			FatherID: 0,
			Name:     "r1",
		},
		nil,
	}
	r2 = &Route{
		&BoxRoute{
			Model: gorm.Model{
				ID: 2,
			},
			FatherID: 1,
			Name:     "r2",
		},
		nil,
	}
	r3 = &Route{
		&BoxRoute{
			Model: gorm.Model{
				ID: 3,
			},
			FatherID: 1,
			Name:     "r3",
		},
		nil,
	}
	r4 = &Route{
		&BoxRoute{
			Model: gorm.Model{
				ID: 4,
			},
			FatherID: 1,
			Name:     "r4",
		},
		nil,
	}
	r5 = &Route{
		&BoxRoute{
			Model: gorm.Model{
				ID: 5,
			},
			FatherID: 1,
			Name:     "r5",
		},
		nil,
	}
)

func Test_getRoutes(t *testing.T) {
	type args struct {
		m   map[uint][]*Route
		fid uint
	}
	tests := []struct {
		name string
		args args
		want []*Route
	}{
		{
			"base test",
			args{
				map[uint][]*Route{
					0: []*Route{r1},
					1: []*Route{r3, r4, r5},
				},
				0,
			},
			[]*Route{r1},
		},
		{
			"base test",
			args{
				map[uint][]*Route{
					0: []*Route{r1},
					1: []*Route{r3, r4, r5},
				},
				1,
			},
			[]*Route{r3, r4, r5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRoutes(tt.args.m, tt.args.fid); !reflect.DeepEqual(got, tt.want) {
				t.Logf("routes = %v\n", []*Route{r1, r2, r3, r4, r5})
				t.Errorf("getRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
