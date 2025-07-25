import 'package:hive/hive.dart';
import 'currency.dart';

part 'settings.g.dart';

@HiveType(typeId: 2)
class Settings {
  @HiveField(0)
  final bool isDarkMode;

  @HiveField(1)
  final bool useSystemTheme;

  @HiveField(2)
  final String defaultCurrency;

  @HiveField(3)
  final String? selectedFamilyId;

  @HiveField(4)
  final bool hasCompletedOnboarding;

  Settings({
    this.isDarkMode = false,
    this.useSystemTheme = true,
    this.defaultCurrency = Currency.defaultCode,
    this.selectedFamilyId,
    this.hasCompletedOnboarding = false,
  });

  // Get the default currency as a Currency enum
  Currency get defaultCurrencyEnum => Currency.fromCode(defaultCurrency);

  Settings copyWith({
    bool? isDarkMode,
    bool? useSystemTheme,
    String? defaultCurrency,
    Currency? defaultCurrencyEnum,
    String? selectedFamilyId,
    bool? hasCompletedOnboarding,
  }) {
    return Settings(
      isDarkMode: isDarkMode ?? this.isDarkMode,
      useSystemTheme: useSystemTheme ?? this.useSystemTheme,
      defaultCurrency:
          defaultCurrencyEnum?.code ?? defaultCurrency ?? this.defaultCurrency,
      selectedFamilyId: selectedFamilyId ?? this.selectedFamilyId,
      hasCompletedOnboarding: hasCompletedOnboarding ?? this.hasCompletedOnboarding,
    );
  }
}
