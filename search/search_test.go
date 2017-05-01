package search

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_checkDate(t *testing.T) {
	type args struct {
		date string
		tm   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Correct date and time",
			args: args{
				date: "2017-01-01",
				tm:   "00:00:00",
			},
			want:    "2017-01-01T00:00:00",
			wantErr: false,
		},
		{
			name: "Incorrect date and correct time",
			args: args{
				date: "2017-32-01",
				tm:   "00:00:00",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Incorrect date and correct time",
			args: args{
				date: "xx-yy-zz",
				tm:   "00:00:00",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Correct date and incorrect time",
			args: args{
				date: "2017-01-01",
				tm:   "00:00:77",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Correct date and incorrect time",
			args: args{
				date: "2017-01-01",
				tm:   "ascd",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkDate(tt.args.date, tt.args.tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSearchTerms(t *testing.T) {
	type args struct {
		r *http.Request
	}

	// Successful scenario
	uA := url.Values{}
	uA.Set("table", "aaaa")
	uA.Set("namespace", "aaaa1")
	uA.Set("from_date", "2017-10-10")
	uA.Set("from_time", "10:10:10")
	uA.Set("to_date", "2017-11-11")
	uA.Set("to_time", "11:11:11")
	rA := &http.Request{
		Form: uA,
	}
	wantA := make(map[string]string)
	wantA["namespace"] = "aaaa1"
	wantA["from_date"] = "2017-10-10T10:10:10"
	wantA["to_date"] = "2017-11-11T11:11:11"

	// Unsuccessful scenarios
	uB := url.Values{}
	uB.Set("table", "")
	uB.Set("namespace", "aaaa1")
	uB.Set("from_date", "2017-10-10")
	uB.Set("from_time", "10:10:10")
	uB.Set("to_date", "2017-11-11")
	uB.Set("to_time", "11:11:11")
	rB := &http.Request{
		Form: uB,
	}

	uC := url.Values{}
	uC.Set("table", "aaaa")
	uC.Set("namespace", "")
	uC.Set("from_date", "2017-10-10")
	uC.Set("from_time", "10:10:10")
	uC.Set("to_date", "2017-11-11")
	uC.Set("to_time", "11:11:11")
	rC := &http.Request{
		Form: uC,
	}

	uD := url.Values{}
	uD.Set("table", "aaaa")
	uD.Set("namespace", "aaaa1")
	uD.Set("from_date", "")
	uD.Set("from_time", "10:10:10")
	uD.Set("to_date", "2017-11-11")
	uD.Set("to_time", "11:11:11")
	rD := &http.Request{
		Form: uD,
	}

	uE := url.Values{}
	uE.Set("table", "aaaa")
	uE.Set("namespace", "aaaa1")
	uE.Set("from_date", "2017-10-10")
	uE.Set("from_time", "")
	uE.Set("to_date", "2017-11-11")
	uE.Set("to_time", "11:11:11")
	rE := &http.Request{
		Form: uE,
	}

	uF := url.Values{}
	uF.Set("table", "aaaa")
	uF.Set("namespace", "aaaa1")
	uF.Set("from_date", "2017-10-10")
	uF.Set("from_time", "10:10:10")
	uF.Set("to_date", "")
	uF.Set("to_time", "11:11:11")
	rF := &http.Request{
		Form: uF,
	}

	uG := url.Values{}
	uG.Set("table", "aaaa")
	uG.Set("namespace", "aaaa1")
	uG.Set("from_date", "2017-10-10")
	uG.Set("from_time", "10:10:10")
	uG.Set("to_date", "2017-11-11")
	uG.Set("to_time", "")
	rG := &http.Request{
		Form: uG,
	}

	tests := []struct {
		name    string
		args    args
		want    string
		want1   map[string]string
		wantErr bool
	}{
		{
			name: "Correct search form input",
			args: args{
				r: rA,
			},
			want:    "aaaa",
			want1:   wantA,
			wantErr: false,
		},
		{
			name: "Incorrect search form input - missing table value",
			args: args{
				r: rB,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect search form input - missing namespace value",
			args: args{
				r: rC,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect search form input - missing from_date value",
			args: args{
				r: rD,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect search form input - missing from_time value",
			args: args{
				r: rE,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect search form input - missing to_date value",
			args: args{
				r: rF,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect search form input - missing to_time value",
			args: args{
				r: rG,
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getSearchTerms(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSearchTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSearchTerms() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getSearchTerms() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
