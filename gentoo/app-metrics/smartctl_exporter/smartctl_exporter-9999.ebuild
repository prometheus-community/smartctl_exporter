# Copyright 2022 The Prometheus Authors
# Copyright 2019 Maxim "Sheridan" Gorlov
# Distributed under the terms of the Apache License Version 2.0

EAPI=6
EGIT_REPO_URI="https://github.com/prometheus-community/${PN}"
inherit git-r3

DESCRIPTION="Exporting S.M.A.R.T. metrics"
HOMEPAGE="https://github.com/prometheus-community/smartctl_exporter"
LICENSE="Apache-2.0"
SLOT="0"
RDEPEND="sys-apps/smartmontools"
DEPEND="dev-lang/go"
KEYWORDS="~amd64 ~ppc ~x86 ~arm"

src_unpack() {
        git-r3_src_unpack
        cd "${S}"
        make get
}

src_compile() {
	make build
}

src_install() {
	newbin bin/${PN}-${PV} ${PN}
	dodoc "README.md"
        insinto /etc
}
