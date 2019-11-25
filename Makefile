tutorial:
	@# todo: have this actually run some kind of tutorial wizard?
	@echo "Please read the 'Makefile' file to go through this tutorial"

var-host:
	export foo=bar
	echo "foo=[$$foo]"

var-kept:
	export foo=bar; \
	echo "foo=[$$foo]"

result.txt: source.txt
	@echo "building result.txt from source.txt"
	cp source.txt result.txt

source.txt:
	@echo "building source.txt"
	echo "this is the source." > source.txt

# $@ The file that is being made right now by ths rule(aka the target).
# $< The input file(that is, the first prerequisite int the list).
# $^ This is the list of ALL input files, not just the first one.
# $? All the input files that are newer than the target.
# $$ Aliteral $ characterinside of the rules section.
# $* The "stem" part that matched in the rule definition's % bit.
# Other special syntax $(@D) and $(@F) to refer to just the dir and file portions of $@, respectively.

result-using-var.txt: source.txt
	@echo "building result-using-var.txt using the $$< and $$@ vars."
	cp $< $@

srcfiles := $(shell echo src/{00..03}.txt)

src/%.txt:
	@# First things first, create the dir if it doesn't exist.
	@# Prepend with @ because srsly who cares about dir creation.
	@[[ -d src ]] || mkdir src
	@# Then, we just echo some data into the file,
	@# The $* expands to the "stem" bit matched by %.
	@# So, we get a bunch of files with numeric names, containing their number.
	echo $* > $@

source: $(srcfiles)

dest/%.txt: src/%.txt
	@[ -d dest ] || mkdir dest
	cp $< $@

destfiles := $(patsubst src/$.txt,dest/%.txt,$(srcfiles))
destination: $(destfiles)

kitty: $(destfiles)
	@# Remember, $< is the input file, but $^ is ALL the input files.
	@# Cat them into the kitty.
	cat $^ > $@

test: kitty
	@echo "miao" && echo "test all pass!"

clean:
	@rm -rf *.txt src dest kitty

badkitty:
	$(MAKE) kitty # The special var $(MAKE) means "the make currently in use"
	false # <-- this will fail
	echo "should not got there."

.PHONY: source destination clean test badkitty