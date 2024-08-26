package wordenc

import (
	"path/filepath"
	"testing"

	"github.com/zsrv/rs-server-225/internal/projectpath"
	"github.com/zsrv/rs-server-225/jagex2/wordpack"
)

func Test_filter(t *testing.T) {
	wf, err := LoadWordFilter(filepath.Join(projectpath.Root, "data", "pack"))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "runescape dot com",
			args: args{input: wordpack.ToSentenceCase("runescape dot com")},
			want: "*****************",
		},
		{
			name: "runescape(dot)com",
			args: args{input: wordpack.ToSentenceCase("runescape(dot)com")},
			want: "*****************",
		},
		{
			name: "runescape.com",
			args: args{input: wordpack.ToSentenceCase("runescape.com")},
			want: "*************",
		},
		{
			name: "well fuck man",
			args: args{input: wordpack.ToSentenceCase("well fuck man")},
			want: "Well **** man",
		},
		{
			name: "google.com\\chrome",
			args: args{input: wordpack.ToSentenceCase("google.com\\chrome")},
			want: "*****************",
		},
		{
			name: "g00gle (dot) c0m \\ chrome",
			args: args{input: wordpack.ToSentenceCase("g00gle (dot) c0m \\ chrome")},
			want: "*************************",
		},
		{
			name: "g00gle hey (d0t) (c0m) (slash) chr0me",
			args: args{input: wordpack.ToSentenceCase("g00gle hey (d0t) (c0m) (slash) chr0me")},
			want: "G00gle ******************************",
		},
		{
			name: "test(at)gmail(dot)com",
			args: args{input: wordpack.ToSentenceCase("test(at)gmail(dot)com")},
			want: "Test(at)*************",
		},
		{
			name: "EF^&N*DGTFbnds7fyt8a^NEAFTBfdasBTFNB(*DS&YT",
			args: args{input: wordpack.ToSentenceCase("EF^&N*DGTFbnds7fyt8a^NEAFTBfdasBTFNB(*DS&YT")},
			want: "Ef^&n*dgtfbnds7fyt8a^neaftbfdasbtfnb(*ds&yt",
		},
		{
			name: "well, anyways. what is up homie? lol fuck i didnt mean to!!!!! ? uhhhh :)",
			args: args{input: wordpack.ToSentenceCase("well, anyways. what is up homie? lol fuck i didnt mean to!!!!! ? uhhhh :)")},
			want: "Well, anyways. What is up homie? lol **** i didnt mean to!!!!! ? Uhhhh :)",
		},
		{
			name: "im so fucking bored",
			args: args{input: wordpack.ToSentenceCase("im so fucking bored")},
			want: "Im so ****ing bored",
		},
		{
			name: "i focking hate this sh!t",
			args: args{input: wordpack.ToSentenceCase("i focking hate this sh!t")},
			want: "I ****ing hate thi******",
		},
		{
			name: "fuckign bit ch",
			args: args{input: wordpack.ToSentenceCase("fuckign bit ch")},
			want: "****ign ******",
		},
		{
			name: "this is a badword.com",
			args: args{input: wordpack.ToSentenceCase("this is a badword.com")},
			want: "This is a ***********",
		},
		{
			name: "this is a badword.org",
			args: args{input: wordpack.ToSentenceCase("this is a badword.org")},
			want: "This is a ***********",
		},
		{
			name: "this is a badword.net",
			args: args{input: wordpack.ToSentenceCase("this is a badword.net")},
			want: "This is a ***********",
		},
		{
			name: "this is a badword.xyz",
			args: args{input: wordpack.ToSentenceCase("this is a badword.xyz")},
			want: "This is a badword.Xyz",
		},
		{
			name: "badword.com is good",
			args: args{input: wordpack.ToSentenceCase("badword.com is good")},
			want: "*********** is good",
		},
		{
			name: "badword.org is good",
			args: args{input: wordpack.ToSentenceCase("badword.org is good")},
			want: "*********** is good",
		},
		{
			name: "badword.net is good",
			args: args{input: wordpack.ToSentenceCase("badword.net is good")},
			want: "*********** is good",
		},
		{
			name: "fragment.xyz is good",
			args: args{input: wordpack.ToSentenceCase("fragment.xyz is good")},
			want: "Fragment.Xyz is good",
		},
		{
			name: "badword (dot) com",
			args: args{input: wordpack.ToSentenceCase("badword (dot) com")},
			want: "*****************",
		},
		{
			name: "badword (d0t) c0m",
			args: args{input: wordpack.ToSentenceCase("badword (d0t) c0m")},
			want: "*****************",
		},
		{
			name: "badword [dot] com",
			args: args{input: wordpack.ToSentenceCase("badword [dot] com")},
			want: "*****************",
		},
		{
			name: "badword [d0t] c0m",
			args: args{input: wordpack.ToSentenceCase("badword [d0t] c0m")},
			want: "*****************",
		},
		{
			name: "badword {dot} com",
			args: args{input: wordpack.ToSentenceCase("badword {dot} com")},
			want: "*****************",
		},
		{
			name: "badword {d0t} c0m",
			args: args{input: wordpack.ToSentenceCase("badword {d0t} c0m")},
			want: "*****************",
		},
		{
			name: "badword slash com",
			args: args{input: wordpack.ToSentenceCase("badword slash com")},
			want: "Badword slash ***",
		},
		{
			name: "badword sl4sh c0m",
			args: args{input: wordpack.ToSentenceCase("badword sl4sh c0m")},
			want: "Badword sl4sh ***",
		},
		{
			name: "badword slash c0m",
			args: args{input: wordpack.ToSentenceCase("badword slash c0m")},
			want: "Badword slash ***",
		},
		{
			name: "com dot badword",
			args: args{input: wordpack.ToSentenceCase("com dot badword")},
			want: "*** dot badword",
		},
		{
			name: "c0m d0t badword",
			args: args{input: wordpack.ToSentenceCase("c0m d0t badword")},
			want: "*** d0t badword",
		},
		{
			name: "com, badword",
			args: args{input: wordpack.ToSentenceCase("com, badword")},
			want: "***, badword",
		},
		{
			name: "c0m, badword",
			args: args{input: wordpack.ToSentenceCase("c0m, badword")},
			want: "***, badword",
		},
		{
			name: "----vv vv vv rswalmart  c - 0 - nn sell cheap gold 1000k \"=\"2.1(.u.s'd)",
			args: args{input: wordpack.ToSentenceCase("----vv vv vv rswalmart  c - 0 - nn sell cheap gold 1000k \"=\"2.1(.u.s'd)")},
			want: "----******** rswalmart c - 0 - nn sell cheap gold 1000k \"=\"2.1(.*****)",
		},
		{
			name: "Web:---4 r s_gold_c\"..0..\"\"'|\\/|\"\"cheap rs gold -20 'm'=18.3'$",
			args: args{input: wordpack.ToSentenceCase("Web:---4 r s_gold_c\"..0..\"\"'|\\/|\"\"cheap rs gold -20 'm'=18.3'$")},
			want: "Web:---4 r s_gold_********************* rs gold -20 'm'=18.3'$",
		},
		{
			name: "Cheap sell gold>google open:___'fzf'__c\"..0..'|\\/|\"\">20m=17.23$",
			args: args{input: wordpack.ToSentenceCase("Cheap sell gold>google open:___'fzf'__c\"..0..'|\\/|\"\">20m=17.23$")},
			want: "Cheap sell gold>google open:___'fzf'__******************=17.23$",
		},
		{
			name: "..:::.4 r s g 0 l d..:::c:::0:::/y\\>>>20\"\"m = 18.3----usd.",
			args: args{input: wordpack.ToSentenceCase("..:::.4 r s g 0 l d..:::c:::0:::/y\\>>>20\"\"m = 18.3----usd.")},
			want: "..:::.4 R s g 0 l ****************\\>>>***** = 18.3----Usd.",
		},

		// should not filter

		{
			name: "runescape",
			args: args{input: wordpack.ToSentenceCase("runescape")},
			want: "Runescape",
		},
		{
			name: "hello@man",
			args: args{input: wordpack.ToSentenceCase("hello@man")},
			want: "Hello@man",
		},
		{
			name: "(dot)",
			args: args{input: wordpack.ToSentenceCase("(dot)")},
			want: "(Dot)",
		},
		{
			name: "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
			args: args{input: wordpack.ToSentenceCase("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")},
			want: "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
		},
		{
			name: "#######################",
			args: args{input: wordpack.ToSentenceCase("#######################")},
			want: "#######################",
		},
		{
			name: "hello world",
			args: args{input: wordpack.ToSentenceCase("hello world")},
			want: "Hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(wf, tt.args.input); got != tt.want {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_filter(b *testing.B) {
	wf, err := LoadWordFilter(filepath.Join(projectpath.Root, "data", "pack"))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter(wf, "badword [dot] com")
	}
}
