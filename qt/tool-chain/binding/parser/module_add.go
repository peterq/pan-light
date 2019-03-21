package parser

import "fmt"

func (m *Module) add() {
	switch m.Project {
	case "QtWebEngine":
		if _, ok := State.ClassMap["QWebEngineCertificateError"]; ok {
			return
		}
		c := &Class{
			Name:   "QWebEngineCertificateError",
			Status: "active",
			Module: m.Project,
			Access: "public",
			Functions: []*Function{
				{
					Name:       "error",
					Fullname:   fmt.Sprintf("%v::%v", "QWebEngineCertificateError", "error"),
					Access:     "public",
					Virtual:    "non",
					Meta:       PLAIN,
					Output:     "Error",
					Parameters: []*Parameter{},
					Signature:  "Error error() const",
				},
				{
					Name:       "errorDescription",
					Fullname:   fmt.Sprintf("%v::%v", "QWebEngineCertificateError", "errorDescription"),
					Access:     "public",
					Virtual:    "non",
					Meta:       PLAIN,
					Output:     "QString",
					Parameters: []*Parameter{},
					Signature:  "QString errorDescription() const",
				},
				{
					Name:       "isOverridable",
					Fullname:   fmt.Sprintf("%v::%v", "QWebEngineCertificateError", "isOverridable"),
					Access:     "public",
					Virtual:    "non",
					Meta:       PLAIN,
					Output:     "bool",
					Parameters: []*Parameter{},
					Signature:  "bool isOverridable() const",
				},
				{
					Name:       "url",
					Fullname:   fmt.Sprintf("%v::%v", "QWebEngineCertificateError", "url"),
					Access:     "public",
					Virtual:    "non",
					Meta:       PLAIN,
					Output:     "QUrl",
					Parameters: []*Parameter{},
					Signature:  "QUrl url() const",
				},
			},
			Enums: []*Enum{{
				Name:     "Error",
				Fullname: "QWebEngineCertificateError::Error",
				Status:   "active",
				Access:   "public",
				Values: []*Value{
					{"SslPinnedKeyNotInCertificateChain", "-150"},
					{"CertificateCommonNameInvalid", "-200"},
					{"CertificateDateInvalid", "-201"},
					{"CertificateAuthorityInvalid", "-202"},
					{"CertificateContainsErrors", "-203"},
					{"CertificateNoRevocationMechanism", "-204"},
					{"CertificateUnableToCheckRevocation", "-205"},
					{"CertificateRevoked", "-206"},
					{"CertificateInvalid", "-207"},
					{"CertificateWeakSignatureAlgorithm", "-208"},
					{"CertificateNonUniqueName", "-210"},
					{"CertificateWeakKey", "-211"},
					{"CertificateNameConstraintViolation", "-212"},
					{"CertificateValidityTooLong", "-213"},
					{"CertificateTransparencyRequired", "-214"},
				}},
			},
		}
		c.register(m)
	}
}
