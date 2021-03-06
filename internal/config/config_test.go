package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func TestGetString(t *testing.T) {
	type args struct {
		key  string
		flag string
	}
	tests := []struct {
		name       string
		keyToSet   string
		valueToSet interface{}
		args       args
		want       string
	}{
		{
			name: "can get the flag when viper is not set",
			args: args{
				key:  "some.key",
				flag: "value",
			},
			want: "value",
		},
		{
			name:       "can get the flag when viper is set",
			keyToSet:   "some.key",
			valueToSet: "thevalue",
			args: args{
				key:  "some.key",
				flag: "",
			},
			want: "thevalue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.keyToSet != "" {
				viper.Set(tt.keyToSet, tt.valueToSet)
			}

			if got := GetString(tt.args.key, tt.args.flag); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_RemoveSite(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "remove a site by its hostname",
			args: args{
				site: "anotherexample.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example.test",
						Webroot:  "web",
					},
					{
						Hostname: "anotherexample.test",
						Webroot:  "web",
					},
					{
						Hostname: "finalexample.test",
						Webroot:  "web",
					},
				},
			},
			want: []Site{
				{
					Hostname: "example.test",
					Webroot:  "web",
				},
				{
					Hostname: "finalexample.test",
					Webroot:  "web",
				},
			},
			wantErr: false,
		},
		{
			name: "sites not in the slice return an error",
			args: args{
				site: "doesnotexist.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "example.test",
						Webroot:  "web",
					},
					{
						Hostname: "anotherexample.test",
						Webroot:  "web",
					},
					{
						Hostname: "finalexample.test",
						Webroot:  "web",
					},
				},
			},
			want: []Site{
				{
					Hostname: "example.test",
					Webroot:  "web",
				},
				{
					Hostname: "anotherexample.test",
					Webroot:  "web",
				},
				{
					Hostname: "finalexample.test",
					Webroot:  "web",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}

			err := c.RemoveSite(tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RemoveSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_AddMount(t *testing.T) {
	// TODO fix this and add the full paths
	t.Skip("update to use full paths")

	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		m Mount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Mount
		wantErr bool
	}{
		{
			name: "adds a new mount and sets a default dest path for period references",
			args: args{
				m: Mount{
					Source: "./testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount and sets a default dest path for non-relative references",
			args: args{
				m: Mount{
					Source: "testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount and sets a default dest path for relative",
			args: args{
				m: Mount{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
					Dest:   "/nitro/sites/test-mount",
				},
			},
			wantErr: false,
		},
		{
			name: "adds a new mount",
			args: args{
				m: Mount{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
					Dest:   "/home/ubuntu/sites",
				},
			},
			want: []Mount{
				{
					Source: "~/go/src/github.com/craftcms/nitro/internal/config/testdata/test-mount",
					Dest:   "/home/ubuntu/sites",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.AddMount(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("AddMount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Mounts, tt.want) {
					t.Errorf("AddMount() got = \n%v, \nwant \n%v", c.Mounts, tt.want)
				}
			}
		})
	}
}

func TestConfig_AddSite(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site Site
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "adds a new site",
			args: args{
				site: Site{
					Hostname: "craftdev",
				},
			},
			want: []Site{
				{
					Hostname: "craftdev",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.AddSite(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("AddSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("AddSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_RemoveMountBySiteWebroot(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		webroot string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Mount
		wantErr bool
	}{
		{
			name: "can remove a mount by a site webroot",
			fields: fields{
				Name: "somename",
				PHP:  "7.4",
				Mounts: []Mount{
					{
						Source: "./testdata/test-mount",
						Dest:   "/nitro/sites/testmount",
					},
					{
						Source: "./testdata/test-mount/remove",
						Dest:   "/nitro/sites/remove",
					},
				},
				Sites: []Site{
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep/web",
					},
					{
						Hostname: "remove.test",
						Webroot:  "/nitro/sites/testmount/remove/web",
					},
				},
			},
			args: args{webroot: "/nitro/sites/remove/web"},
			want: []Mount{
				{
					Source: "./testdata/test-mount",
					Dest:   "/nitro/sites/testmount",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RemoveMountBySiteWebroot(tt.args.webroot); (err != nil) != tt.wantErr {
				t.Errorf("RemoveMountBySiteWebroot() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Mounts, tt.want) {
					t.Errorf("RemoveMountBySiteWebroot() got = \n%v, \nwant \n%v", c.Mounts, tt.want)
				}
			}
		})
	}
}

func TestConfig_RemoveSite1(t *testing.T) {
	type fields struct {
		Name      string
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "can remove a site by its hostname",
			fields: fields{
				Name:   "somename",
				PHP:    "7.4",
				CPUs:   "3",
				Disk:   "20G",
				Memory: "4G",
				Sites: []Site{
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep/web",
					},
					{
						Hostname: "remove.test",
						Webroot:  "/nitro/sites/remove/web",
					},
				},
			},
			args: args{hostname: "remove.test"},
			want: []Site{
				{
					Hostname: "keep.test",
					Webroot:  "/nitro/sites/keep/web",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RemoveSite(tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("RemoveSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RemoveSite() got = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_RenameSite(t *testing.T) {
	type fields struct {
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site     Site
		hostname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Site
		wantErr bool
	}{
		{
			name: "remove a site my hostname",
			args: args{
				site: Site{
					Hostname: "old.test",
					Webroot:  "/nitro/sites/old.test",
				},
				hostname: "new.test",
			},
			fields: fields{
				Sites: []Site{
					{
						Hostname: "old.test",
						Webroot:  "/nitro/sites/old.test",
					},
					{
						Hostname: "keep.test",
						Webroot:  "/nitro/sites/keep.test",
					},
				},
			},
			want: []Site{
				{
					Hostname: "new.test",
					Webroot:  "/nitro/sites/new.test",
				},
				{
					Hostname: "keep.test",
					Webroot:  "/nitro/sites/keep.test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if err := c.RenameSite(tt.args.site, tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("RenameSite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(c.Sites, tt.want) {
					t.Errorf("RenameSite() got sites = \n%v, \nwant \n%v", c.Sites, tt.want)
				}
			}
		})
	}
}

func TestConfig_MountExists(t *testing.T) {
	type fields struct {
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		dest string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "existing mounts return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
			},
			args: args{dest: "/nitro/sites/example-site"},
			want: true,
		},
		{
			name: "non-existing mounts return false",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites/example-site",
					},
				},
			},
			args: args{dest: "/nitro/sites/nonexistent-site"},
			want: false,
		},
		{
			name: "parent mounts return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: "./testdata/test-mount",
						Dest:   "/nitro/sites",
					},
					{
						Source: "./testdata/existing-mount",
						Dest:   "/nitro/sites",
					},
				},
			},
			args: args{dest: "/nitro/sites/nonexistent-site"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.MountExists(tt.args.dest); got != tt.want {
				t.Errorf("MountExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SiteExists(t *testing.T) {
	type fields struct {
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		site Site
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "exact sites return true",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "iexist.test",
						Webroot:  "/nitro/sites/iexist.test",
					},
				},
			},
			args: args{site: Site{
				Hostname: "iexist.test",
				Webroot:  "/nitro/sites/iexist.test",
			}},
			want: true,
		},
		{
			name: "exact sites return false",
			fields: fields{
				Sites: []Site{
					{
						Hostname: "iexist.test",
						Webroot:  "/nitro/sites/iexist.test",
					},
				},
			},
			args: args{site: Site{
				Hostname: "idontexist.test",
				Webroot:  "/nitro/sites/idontexist.test",
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.SiteExists(tt.args.site); got != tt.want {
				t.Errorf("SiteExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_DatabaseExists(t *testing.T) {
	type fields struct {
		PHP       string
		CPUs      string
		Disk      string
		Memory    string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		database Database
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "can find an existing database",
			fields: fields{
				Databases: []Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
				},
			},
			args: args{database: Database{
				Engine:  "mysql",
				Version: "5.8",
				Port:    "3306",
			}},
			want: false,
		},
		{
			name: "non-existing databases return false",
			fields: fields{
				Databases: []Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
				},
			},
			args: args{database: Database{
				Engine:  "mysql",
				Version: "5.7",
				Port:    "3306",
			}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.DatabaseExists(tt.args.database); got != tt.want {
				t.Errorf("DatabaseExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_FindMountBySiteWebroot(t *testing.T) {
	type fields struct {
		PHP       string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		webroot string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Mount
	}{
		{
			name: "",
			fields: fields{
				PHP: "7.4",
				Mounts: []Mount{
					{
						Source: "./testdata/dev/nitro",
						Dest:   "/home/ubuntu/sites/nitro",
					},
				},
				Databases: nil,
				Sites: []Site{
					{
						Hostname: "nitro",
						Webroot:  "/home/ubuntu/sites/nitro/www",
					},
				},
			},
			args: args{webroot: "/home/ubuntu/sites/nitro/www"},
			want: &Mount{
				Source: "./testdata/dev/nitro",
				Dest:   "/home/ubuntu/sites/nitro",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.FindMountBySiteWebroot(tt.args.webroot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMountBySiteWebroot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_AlreadyMounted(t *testing.T) {
	_, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		PHP       string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		m Mount
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantMount Mount
	}{
		{
			name: "mounts that have a parent path return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: current,
						Dest:   "/home/ubuntu/sites/example",
					},
				},
			},
			args: args{m: Mount{
				Source: current + "/testdata/new-mount",
				Dest:   "/home/ubuntu/sites/example",
			}},
			want: true,
			wantMount: Mount{
				Source: current,
				Dest:   "/home/ubuntu/sites/example",
			},
		},
		{
			name: "mounts that exists return true",
			fields: fields{
				Mounts: []Mount{
					{
						Source: current,
						Dest:   "/home/ubuntu/sites/example",
					},
				},
			},
			args: args{m: Mount{
				Source: current,
				Dest:   "/home/ubuntu/sites/example",
			}},
			want: true,
			wantMount: Mount{
				Source: current,
				Dest:   "/home/ubuntu/sites/example",
			},
		},
		{
			name:   "mounts that do not exist return false",
			fields: fields{Mounts: nil},
			args: args{m: Mount{
				Source: current,
				Dest:   "/home/ubuntu/sites/example.site",
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			got, m := c.AlreadyMounted(tt.args.m)

			if got != tt.want {
				t.Errorf("AlreadyMounted() = \n%v, \nwant \n%v", got, tt.want)
			}

			if m != tt.wantMount {
				t.Errorf("AlreadyMounted() = \n%v, \nwant \n%v", m, tt.wantMount)
			}
		})
	}
}

func TestConfig_DatabaseEnginesAsList(t *testing.T) {
	dbs := []Database{
		{
			Engine:  "mysql",
			Version: "5.7",
			Port:    "3306",
		},
		{
			Engine:  "mysql",
			Version: "5.6",
			Port:    "33061",
		},
		{
			Engine:  "postgres",
			Version: "12",
			Port:    "5432",
		},
	}

	type fields struct {
		PHP       string
		Mounts    []Mount
		Databases []Database
		Sites     []Site
	}
	type args struct {
		engine string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "Can get a list of engines limited by the type of mysql",
			fields: fields{Databases: dbs},
			args:   args{engine: "mysql"},
			want:   []string{"mysql_5.7_3306", "mysql_5.6_33061"},
		},
		{
			name:   "Can get a list of engines limited by the type postgres",
			fields: fields{Databases: dbs},
			args:   args{engine: "postgres"},
			want:   []string{"postgres_12_5432"},
		},
		{
			name:   "All databases are returned when the engine is not provided",
			fields: fields{Databases: dbs},
			want:   []string{"mysql_5.7_3306", "mysql_5.6_33061", "postgres_12_5432"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PHP:       tt.fields.PHP,
				Mounts:    tt.fields.Mounts,
				Databases: tt.fields.Databases,
				Sites:     tt.fields.Sites,
			}
			if got := c.DatabaseEnginesAsList(tt.args.engine); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseEnginesAsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
