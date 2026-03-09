# Evolipia Radar - Mobile App

Flutter app untuk baca news dari Evolipia Radar.

## Setup

### 1. Install Flutter

Download: https://flutter.dev/docs/get-started/install

### 2. Clone & Install Dependencies

```bash
git clone https://github.com/hidatara-ds/evolipia-radar.git
cd evolipia-radar
git checkout mobile-app
flutter pub get
```

### 3. Configure Database

Edit `lib/config.dart`:

```dart
class DatabaseConfig {
  static const String host = 'ep-xxx.us-east-2.aws.neon.tech';
  static const String database = 'neondb';
  static const String username = 'your_username';
  static const String password = 'your_password';
}
```

Get credentials from Neon dashboard.

### 4. Run on Phone

```bash
flutter run
```

### 5. Build APK

```bash
flutter build apk --release
```

APK location: `build/app/outputs/flutter-apk/app-release.apk`

Transfer to phone and install!

## Features

- ✅ News list (last 7 days)
- ✅ Sort by AI score
- ✅ Category badges
- ✅ Pull to refresh
- ✅ Open in browser
- ✅ Material Design 3

## Stack

- Flutter 3.0+
- PostgreSQL (via Neon)
- Material Design 3
