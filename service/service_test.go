package service

//func TestService_getCTTConfig(t *testing.T) {
//	type fields struct {
//		TunnelerCommand command_tunneler.TunnelerCommand
//		CommandRunner   command_runner.CommandRunner
//		Log             *logrus.Logger
//		SshConfig       *config.SshConfig
//		commandLineArgs []string
//		listOfCommands  []string
//	}
//	type args struct {
//		cttArgs []string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *mytypes.CTTConfig
//		wantErr bool
//	}{
//		{
//			name:   "happy-flow",
//			fields: fields{},
//			args: args{
//				cttArgs: []string{"-t", "test"},
//			},
//			want: &mytypes.CTTConfig{
//				TunnelConfigName: "test",
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Service{}
//			got, err := s.getCTTConfig(tt.args.cttArgs)
//			assert.True(t, utils.AssertError(tt.wantErr, err), fmt.Sprintf("wantErr: %t, got %s", tt.wantErr, err))
//			assert.Equal(t, tt.want, got)
//		})
//	}
//}
