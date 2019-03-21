package parser

func sailfishModule() *Module {
	return &Module{
		Project: "QtSailfish",
		Namespace: &Namespace{
			Classes: []*Class{
				{
					Name:   "SailfishApp",
					Access: "public",
					Module: "QtSailfish",
					Functions: []*Function{ //TODO: should be in Namespace.Functions
						{
							Name:      "application",
							Fullname:  "SailfishApp::application",
							Access:    "public",
							Virtual:   "non",
							Meta:      PLAIN,
							Static:    true,
							Output:    "QGuiApplication*",
							Signature: "()",
							Parameters: []*Parameter{
								{
									Name:  "argc",
									Value: "int &",
								},
								{
									Name:  "argv",
									Value: "char **",
								},
							},
						},
						{
							Name:      "main",
							Fullname:  "SailfishApp::main",
							Access:    "public",
							Virtual:   "non",
							Meta:      PLAIN,
							Static:    true,
							Output:    "int",
							Signature: "()",
							Parameters: []*Parameter{
								{
									Name:  "argc",
									Value: "int &",
								},
								{
									Name:  "argv",
									Value: "char **",
								},
							},
						},
						{
							Name:      "createView",
							Fullname:  "SailfishApp::createView",
							Access:    "public",
							Virtual:   "non",
							Meta:      PLAIN,
							Static:    true,
							Output:    "QQuickView*",
							Signature: "()",
						},
						{
							Name:      "pathTo",
							Fullname:  "SailfishApp::pathTo",
							Access:    "public",
							Virtual:   "non",
							Meta:      PLAIN,
							Static:    true,
							Output:    "QUrl",
							Signature: "pathTo(const QString &filename)",
							Parameters: []*Parameter{
								{
									Name:  "filename",
									Value: "QString &",
								},
							},
						},
						{
							Name:      "pathToMainQml",
							Fullname:  "SailfishApp::pathToMainQml",
							Access:    "public",
							Virtual:   "non",
							Meta:      PLAIN,
							Static:    true,
							Output:    "QUrl",
							Signature: "pathToMainQml()",
						},
					},
				},
			},
		},
	}
}
