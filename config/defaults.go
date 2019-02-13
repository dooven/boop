package config

var CONFIG_DEFAULTS = Config{
	Users: []string{
		"<replace-this-user>",
	},
	Regions: []RegionOption{
		{
			Region: "ap-northeast-1",
			Name:   "Asia Pacific (Tokyo)",
		},
		{
			Region: "ap-northeast-2",
			Name:   "Asia Pacific (Seoul)",
		},
		{
			Region: "jap-south-1",
			Name:   "Asia Pacific (Mumbai)",
		},
		{
			Region: "ap-southeast-1",
			Name:   "Asia Pacific (Singapore)",
		},
		{
			Region: "ap-southeast-2",
			Name:   "Asia Pacific (Sydney)",
		},
		{
			Region: "ca-central-1",
			Name:   "Canada (Central)",
		},
		{
			Region: "eu-central-1",
			Name:   "EU (Frankfurt)",
		},
		{
			Region: "eu-north-1",
			Name:   "EU (Stockholm)",
		},
		{
			Region: "eu-west-1",
			Name:   "EU (Ireland)",
		},
		{
			Region: "eu-west-2",
			Name:   "EU (London)",
		},
		{
			Region: "us-east-1",
			Name:   "US East (N. Virginia)",
		},
		{
			Region: "us-east-2",
			Name:   "US East (Ohio)",
		},
		{
			Region: "us-west-1",
			Name:   "US West (N. California)",
		},
		{
			Region: "us-west-2",
			Name:   "US West (Oregon)",
		},
	},
}
