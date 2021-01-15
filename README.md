# SD Card Copy script

SD Card Copy script copies assets (images, videos or any file) from the SD Card and upload the difference to an AWS S3 Bucket by checking which assets are not present there. It was created as a backup mechanism as the S3 bucket could be configured to use different [Storage classes](https://aws.amazon.com/s3/storage-classes/), like Glacier, to decrease the cost. 

## Configuration

The scripts read an `assets.yaml` configuration file that must exist where the script is executed as this file describes how the script will copy the assets to AWS S3.

### Example:
```yaml
- description: "JPG images"
  sdCardRegex: "*/DCIM/*.JPG"
  s3BucketName: "assets-backup"
  s3BucketPrefix: "Media/Pictures/JPG"
- description: "RAW images"
  sdCardRegex: "*/DCIM/*.ARW"
  s3BucketName: "assets-backup"
  s3BucketPrefix: "Media/Pictures/RAW"
- description: "Video assets"
  sdCardRegex: "*/PRIVATE/*/CLIP/*"
  s3BucketName: "assets-backup"
  s3BucketPrefix: "Media/Videos"
```

## Notes
- It assumes that the environment has the awscli configured
- It requires the command `find` and `mount`
- It assumes that the SD Card uses exfat as a type and only one SD Card is mounted.
