JAVA_SRC := $(shell find src -name '*.java')
OUT_DIR := out/classes
JAR := target/assinador.jar
POLICY := https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2

.PHONY: all clean test samples

all: $(JAR)

$(JAR): $(JAVA_SRC)
	mkdir -p $(OUT_DIR) target
	javac --release 21 -encoding UTF-8 -d $(OUT_DIR) $(JAVA_SRC)
	jar --create --file $(JAR) --main-class br.ufg.hubsaude.assinador.Main -C $(OUT_DIR) .

samples:
	mkdir -p target
	@UUID=11111111-1111-4111-8111-111111111111; \
	printf '%s\n' '{"resourceType":"Bundle","type":"collection","entry":[{"fullUrl":"urn:uuid:'"$$UUID"'","resource":{"resourceType":"Patient","id":"paciente-1","name":[{"text":"Paciente Teste"}]}}]}' > target/bundle.json; \
	printf '%s\n' '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:'"$$UUID"'"}],"recorded":"2026-01-01T00:00:00Z"}' > target/provenance.json; \
	printf '%s\n' '{"type":"PEM","simulation":"true"}' > target/crypto.json; \
	printf '%s\n' '{"verification":{"maxRetries":3,"ocspTimeout":20,"crlTimeout":20}}' > target/config.json
	@rm -f target/test-keystore.p12 target/test-cert.pem
	keytool -genkeypair -alias runner -keyalg RSA -keysize 2048 -validity 3650 -dname "CN=Runner Test,O=UFG,C=BR" -keystore target/test-keystore.p12 -storetype PKCS12 -storepass changeit -keypass changeit >/dev/null 2>&1
	keytool -exportcert -alias runner -keystore target/test-keystore.p12 -storepass changeit -rfc -file target/test-cert.pem >/dev/null 2>&1
	@printf '["' > target/certs.json
	@grep -v -- '-----' target/test-cert.pem | tr -d '\r\n' >> target/certs.json
	@printf '"]\n' >> target/certs.json
	@date +%s > target/timestamp.txt

test: $(JAR) samples
	@TS=$$(cat target/timestamp.txt); \
	java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance.json --crypto-material target/crypto.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json --signer 'Teste Runner' > target/assinatura.json; \
	java -jar $(JAR) validate --signature target/assinatura.json --timestamp $$TS --policy '$(POLICY)' --bundle target/bundle.json --provenance target/provenance.json --config target/config.json > target/validacao.json; \
	grep -q '"valid":true' target/validacao.json
	@TS=$$(cat target/timestamp.txt); \
	printf '%s\n' '{"type":"TOKEN","pin":"123456","identifier":"runner","simulation":"true"}' > target/crypto-token.json; \
	java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance.json --crypto-material target/crypto-token.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json --signer 'Teste Runner Token' > target/assinatura-token.json; \
	grep -q '"resourceType":"Signature"' target/assinatura-token.json
	@TS=$$(cat target/timestamp.txt); \
	printf '%s\n' '{"type":"TOKEN","identifier":"runner"}' > target/crypto-token-sem-pin.json; \
	! java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance.json --crypto-material target/crypto-token-sem-pin.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json > target/token-sem-pin.out 2> target/token-sem-pin.err; \
	grep -q 'CRYPTO.PIN-REQUIRED' target/token-sem-pin.err
	@TS=$$(cat target/timestamp.txt); \
	printf '%s\n' '{"type":"TOKEN","pin":"123456","identifier":"runner","pkcs11Library":"/arquivo/pkcs11-inexistente.so"}' > target/crypto-token-lib-inexistente.json; \
	! java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance.json --crypto-material target/crypto-token-lib-inexistente.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json > target/token-lib-inexistente.out 2> target/token-lib-inexistente.err; \
	grep -q 'PKCS11.LIBRARY-NOT-FOUND' target/token-lib-inexistente.err
	@TS=$$(cat target/timestamp.txt); \
	! java -jar $(JAR) sign --provenance target/provenance.json --crypto-material target/crypto.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json > target/missing-bundle.out 2> target/missing-bundle.err; \
	grep -q 'INPUT.MISSING-PARAMETER' target/missing-bundle.err
	@TS=$$(cat target/timestamp.txt); \
	! java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance.json --crypto-material target/crypto.json --cert-chain target/certs.json --timestamp $$TS --policy 'politica-invalida' --config target/config.json > target/invalid-policy.out 2> target/invalid-policy.err; \
	grep -q 'POLICY.ID-INVALID' target/invalid-policy.err
	@TS=$$(cat target/timestamp.txt); \
	printf '%s\n' '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:22222222-2222-4222-8222-222222222222"}]}' > target/provenance-target-invalido.json; \
	! java -jar $(JAR) sign --bundle target/bundle.json --provenance target/provenance-target-invalido.json --crypto-material target/crypto.json --cert-chain target/certs.json --timestamp $$TS --policy '$(POLICY)' --config target/config.json > target/invalid-provenance.out 2> target/invalid-provenance.err; \
	grep -q 'FORMAT.PROVENANCE-TARGET-NOT-FOUND' target/invalid-provenance.err
	@TS=$$(cat target/timestamp.txt); \
	printf '%s\n' '{"resourceType":"OperationOutcome"}' > target/not-signature.json; \
	! java -jar $(JAR) validate --signature target/not-signature.json --timestamp $$TS --policy '$(POLICY)' --bundle target/bundle.json --provenance target/provenance.json --config target/config.json > target/not-signature.out 2> target/not-signature.err; \
	grep -q 'FORMAT.SIGNATURE-DATA-INVALID' target/not-signature.err
	@PORT=18080; TS=$$(cat target/timestamp.txt); \
	java -jar $(JAR) server --port $$PORT --timeout 1 > target/server.out 2> target/server.err & PID=$$!; \
	trap 'curl -fsS -X POST http://127.0.0.1:'$$PORT'/shutdown >/dev/null 2>&1 || true; kill '$$PID' >/dev/null 2>&1 || true' EXIT; \
	for i in $$(seq 1 30); do curl -fsS http://127.0.0.1:$$PORT/health >/dev/null 2>&1 && break; sleep 0.2; done; \
	curl -fsS http://127.0.0.1:$$PORT/health >/dev/null; \
	printf '{"bundle":"target/bundle.json","provenance":"target/provenance.json","cryptoMaterial":"target/crypto.json","certificateChain":"target/certs.json","timestamp":"%s","policy":"%s","config":"target/config.json","signer":"Teste Runner HTTP"}' $$TS '$(POLICY)' > target/http-sign-body.json; \
	curl -fsS -X POST -H 'Content-Type: application/json' --data @target/http-sign-body.json http://127.0.0.1:$$PORT/sign > target/http-assinatura.json; \
	grep -q '"resourceType":"Signature"' target/http-assinatura.json; \
	printf '{"signature":"target/http-assinatura.json","timestamp":"%s","policy":"%s","bundle":"target/bundle.json","provenance":"target/provenance.json","config":"target/config.json"}' $$TS '$(POLICY)' > target/http-validate-body.json; \
	curl -fsS -X POST -H 'Content-Type: application/json' --data @target/http-validate-body.json http://127.0.0.1:$$PORT/validate > target/http-validacao.json; \
	grep -q '"valid":true' target/http-validacao.json; \
	curl -fsS -X POST http://127.0.0.1:$$PORT/shutdown >/dev/null; \
	wait $$PID || true; \
	trap - EXIT
	@echo 'Testes locais do assinador.jar passaram.'

clean:
	rm -rf out target
