# MySite

A website outlining me, myself, and I.

# Dependencies

- `pandoc`

# Blog

## Hashing

Content to hash:
    1. Title
    2. Author(s) (separated by newline, '\n'; no trailing)
    3. Category(s) (separated by newline, '\n'; no trailing)
    4. Big-endian 64-bit timestamp (unsigned)
    5. Big-endian 32-bit timezone offset
    6. Previous hash hex (or 32 zero bytes, i.e., hex of 64 zeros)
If, on the REMOTE chance, a rehash is required, increment the timestamp.

# Notes

## App Store Icons

- App Store link icons downloaded: 5/21/2025
    - Guidelines: https://developer.apple.com/app-store/marketing/guidelines/#section-badges
- Google Play link icons downloaded: 5/21/2025
    - Guidelines: https://partnermarketinghub.withgoogle.com/brands/google-play/visual-identity/badge-guidelines/

# TODO
- [ ] Add static page from site below
    - https://support.google.com/admob/answer/9363762?hl=en#zippy=%2Cview-an-example-app-adstxt-file
- [ ] Possibly make the IPs environ vars
- [ ] Fix indicators files
- [ ] When submitting product review, remove platform options invalid for product
- [ ] Admin
- [ ] Log when product/issue created/edited
- [ ] Option to pass config file with CLI opts
- [ ] Request tracking/logging
- [ ] Make robots.txt requests to APIs to add to overall robots.txt
- [ ] Merge instances of "style" and "script" into "head" block in templates


# todo
- [ ] Auto-generate Copyright year
