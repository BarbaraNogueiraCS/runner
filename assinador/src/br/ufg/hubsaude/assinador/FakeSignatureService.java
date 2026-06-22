package br.ufg.hubsaude.assinador;

/**
 * Implementação simulada da interface SignatureService.
 *
 * Esta classe concentra o comportamento de assinatura/validação fake exigido na Sprint 2.
 * A implementação reaproveita GuideSignatureService, que monta respostas simuladas com
 * estrutura FHIR Signature/JWS e valida os parâmetros antes de processar a operação.
 */
final class FakeSignatureService extends GuideSignatureService {}
