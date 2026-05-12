package br.ufg.hubsaude.assinador.pkcs11;

public class PKCS11Adapter {
    public boolean checkAvailability() {
        // Ponto de extensão para integração com SunPKCS11, SoftHSM2, token USB ou smart card.
        // A assinatura real não faz parte do escopo desta versão acadêmica.
        return false;
    }
}
