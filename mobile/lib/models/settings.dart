import 'package:hive/hive.dart';

part 'settings.g.dart';

@HiveType(typeId: 2)
class Settings {
  @HiveField(0)
  final bool isDarkMode;

  @HiveField(1)
  final bool useSystemTheme;

  @HiveField(2)
  final String defaultCurrency;

  Settings({
    this.isDarkMode = false,
    this.useSystemTheme = true,
    this.defaultCurrency = 'USD',
  });

  Settings copyWith({
    bool? isDarkMode,
    bool? useSystemTheme,
    String? defaultCurrency,
  }) {
    return Settings(
      isDarkMode: isDarkMode ?? this.isDarkMode,
      useSystemTheme: useSystemTheme ?? this.useSystemTheme,
      defaultCurrency: defaultCurrency ?? this.defaultCurrency,
    );
  }
}