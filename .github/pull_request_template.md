# Pull Request

## Description

Brief description of the changes in this PR.

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Performance improvement

## Related Issues

Closes #[issue number]

## Changes Made

- [ ] Change 1
- [ ] Change 2
- [ ] Change 3

## Testing

### Testing Done

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] All existing tests pass

### Email Testing

- [ ] SMTP connection tested
- [ ] Email delivery verified
- [ ] Template rendering checked
- [ ] Error handling validated

### Testing Instructions

Provide step-by-step instructions for reviewers to test your changes:

1. Set up SMTP server or use test credentials
2. Configure environment: `cp .env.example .env`
3. Start server: `go run ./cmd/server`
4. Send test email: `curl -X POST http://localhost:8080/send ...`

## Code Quality

- [ ] Code follows Go best practices and project style guidelines
- [ ] Self-review of code completed
- [ ] Code is properly documented with Go comments
- [ ] No debug code left in (fmt.Print, log statements, etc.)
- [ ] Error handling is comprehensive
- [ ] Context is properly handled for cancellation

## Security Considerations

- [ ] No sensitive data (passwords, API keys) in code
- [ ] SMTP credentials handled securely
- [ ] Input validation implemented
- [ ] Rate limiting considered

## Breaking Changes

- [ ] No breaking changes
- [ ] Breaking changes documented in CHANGELOG.md
- [ ] Migration guide provided (if applicable)

## Performance Impact

- [ ] No performance impact
- [ ] Performance improvement (include benchmarks)
- [ ] Potential performance regression (documented and justified)

## Reviewers

@mention specific people you want to review this PR

## Additional Notes

Any additional information that would be helpful for reviewers.