JAVA_SRC := $(shell find src -name '*.java')
OUT_DIR := out/classes
JAR := target/assinador.jar

.PHONY: all clean test

all: $(JAR)

$(JAR): $(JAVA_SRC)
	mkdir -p $(OUT_DIR) target
	javac --release 21 -encoding UTF-8 -d $(OUT_DIR) $(JAVA_SRC)
	jar --create --file $(JAR) --main-class br.ufg.hubsaude.assinador.Main -C $(OUT_DIR) .

test: $(JAR)
	@echo 'documento de teste' > target/documento.txt
	@java -jar $(JAR) sign --input target/documento.txt --signer 'Teste Runner' > target/assinatura.json
	@java -jar $(JAR) validate --signature target/assinatura.json > target/validacao.json
	@grep -q '"valid":true' target/validacao.json
	@echo 'Testes locais do assinador.jar passaram.'

clean:
	rm -rf out target
