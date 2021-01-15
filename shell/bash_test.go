package shell

import "testing"

func TestBash_Execute(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "It should output the echo command",
			args:    args{
				cmd: "echo test123",
			},
			want:    "test123\n",
			wantErr: false,
		},
		{
			name:    "It should fail when the command does not exist",
			args:    args{
				cmd: "test123",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bash{}
			got, err := b.Execute(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
