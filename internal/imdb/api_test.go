package imdb

import "testing"

func TestIsRatingsHyperlink(t *testing.T) {
	tests := []struct {
		name string
		href string
		want bool
	}{
		{
			name: "legacy ur-prefixed ratings export href",
			href: "/user/ur36148206/ratings/?ref_=exp_t_2",
			want: true,
		},
		{
			name: "new p-prefixed ratings export href",
			href: "/user/p.4ngm3d4h2yvdtd2k7pptri4nvm/ratings/?ref_=exp_t_2",
			want: true,
		},
		{
			name: "list export href is not a ratings href",
			href: "/list/ls050054608/?ref_=exp_t_1",
			want: false,
		},
		{
			name: "watchlist user href is not a ratings href",
			href: "/user/ur36148206/watchlist/",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRatingsHyperlink(tt.href); got != tt.want {
				t.Fatalf("isRatingsHyperlink(%q) = %v, want %v", tt.href, got, tt.want)
			}
		})
	}
}

func TestUserIDExtract(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "new p-prefixed user id from watchlist redirect",
			rawURL: "https://www.imdb.com/user/p.4ngm3d4h2yvdtd2k7pptri4nvm/watchlist/",
			want:   "p.4ngm3d4h2yvdtd2k7pptri4nvm",
		},
		{
			name:   "legacy ur-prefixed user id",
			rawURL: "https://www.imdb.com/user/ur12345678/watchlist/",
			want:   "ur12345678",
		},
		{
			name:   "localised url with language prefix",
			rawURL: "https://www.imdb.com/de/user/p.4ngm3d4h2yvdtd2k7pptri4nvm/watchlist/",
			want:   "p.4ngm3d4h2yvdtd2k7pptri4nvm",
		},
		{
			name:   "url with query string",
			rawURL: "https://www.imdb.com/user/p.abc/ratings/?ref_=wl_sm",
			want:   "p.abc",
		},
		{
			name:    "url without a user segment",
			rawURL:  "https://www.imdb.com/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userIDExtract(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userIDExtract(%q) error = %v, wantErr %v", tt.rawURL, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("userIDExtract(%q) = %q, want %q", tt.rawURL, got, tt.want)
			}
		})
	}
}
