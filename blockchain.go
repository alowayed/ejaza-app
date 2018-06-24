package main

func saveCertToBlockchain(cert Cert) string {
	tierionId := saveCertToTierion(cert)

	return tierionId
}