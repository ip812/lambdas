bootstrap:
	@cp -r hello "$(LAMBDA)"
	@find "$(LAMBDA)" -type f -exec sed -i '' "s/hello/$(LAMBDA)/g" {} +

