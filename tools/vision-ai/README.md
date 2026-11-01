# Vision AI tools

## `image_audit.py`

Inspect dimensions, mode and EXIF presence; optionally write a copy without
EXIF/GPS metadata.

```bash
pip install -r tools/vision-ai/requirements.txt
python tools/vision-ai/image_audit.py photo.jpg
python tools/vision-ai/image_audit.py photo.jpg --strip /tmp/photo-clean.jpg
```

The tool does not display GPS values and does not upload images. It does not
detect objects or identify people.

Before publishing any image: check consent, license, faces, plates and private
installations. `NOT_SURE` is a valid annotation.
