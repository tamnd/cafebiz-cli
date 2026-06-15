package cafebiz

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "cafebiz" {
		t.Errorf("Scheme = %q, want cafebiz", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "cafebiz" {
		t.Errorf("Binary = %q, want cafebiz", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		in      string
		wantTyp string
		wantID  string
		wantErr bool
	}{
		{"https://cafebiz.vn/doanh-nhan-nao-do-190001.chn", "article", "190001", false},
		{"1234567", "article", "1234567", false},
		{"doanh-nhan", "category", "doanh-nhan", false},
		{"startup", "category", "startup", false},
		{"home", "category", "home", false},
		{"", "", "", true},
		{"not-a-category-or-id", "", "", true},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Classify(%q): want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("Classify(%q): %v", tc.in, err)
			continue
		}
		if typ != tc.wantTyp || id != tc.wantID {
			t.Errorf("Classify(%q) = (%q,%q), want (%q,%q)", tc.in, typ, id, tc.wantTyp, tc.wantID)
		}
	}
}

func TestLocate(t *testing.T) {
	cases := []struct {
		typ, id, want string
		wantErr       bool
	}{
		{"article", "190001", "https://cafebiz.vn/article-190001.chn", false},
		{"category", "startup", "https://cafebiz.vn/startup.chn", false},
		{"unknown", "x", "", true},
	}
	for _, tc := range cases {
		got, err := Domain{}.Locate(tc.typ, tc.id)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Locate(%q,%q): want error", tc.typ, tc.id)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("Locate(%q,%q) = (%q,%v), want (%q,nil)", tc.typ, tc.id, got, err, tc.want)
		}
	}
}

func TestHostWiring(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}

	a := &Article{ID: "190001", URL: "https://cafebiz.vn/doanh-nhan-nao-do-190001.chn", Category: "doanh-nhan"}
	u, err := h.Mint(a)
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	if want := "cafebiz://article/190001"; u.String() != want {
		t.Errorf("Mint = %q, want %q", u.String(), want)
	}
}
