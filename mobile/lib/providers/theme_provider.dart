import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:google_fonts/google_fonts.dart';
import '../repositories/settings_repository.dart';

class ThemeProvider with ChangeNotifier {
  bool _isDarkMode = false;
  bool _useSystemTheme = true;
  final SettingsRepository? _settingsRepository;

  ThemeProvider({SettingsRepository? settingsRepository})
    : _settingsRepository = settingsRepository {
    // Load settings if repository is provided
    if (_settingsRepository != null) {
      _loadSettings();
    } else {
      // Initialize with system theme if no repository
      _initializeWithSystemTheme();
    }
  }

  Future<void> _loadSettings() async {
    final settings = _settingsRepository!.getSettings();
    _isDarkMode = settings.isDarkMode;
    _useSystemTheme = settings.useSystemTheme;

    // If using system theme, update dark mode based on system
    if (_useSystemTheme) {
      _initializeWithSystemTheme();
    }

    notifyListeners();
  }

  void _initializeWithSystemTheme() {
    var brightness =
        SchedulerBinding.instance.platformDispatcher.platformBrightness;
    _isDarkMode = brightness == Brightness.dark;
  }

  bool get isDarkMode => _isDarkMode;

  bool get useSystemTheme => _useSystemTheme;

  ThemeMode get themeMode {
    if (_useSystemTheme) {
      return ThemeMode.system;
    }
    return _isDarkMode ? ThemeMode.dark : ThemeMode.light;
  }

  Future<void> toggleTheme() async {
    _useSystemTheme = false;
    _isDarkMode = !_isDarkMode;

    // Persist settings if repository is available
    if (_settingsRepository != null) {
      await _settingsRepository.updateThemeSettings(
        isDarkMode: _isDarkMode,
        useSystemTheme: _useSystemTheme,
      );
    }

    notifyListeners();
  }

  Future<void> setSystemTheme() async {
    _useSystemTheme = true;
    _initializeWithSystemTheme();

    // Persist settings if repository is available
    if (_settingsRepository != null) {
      await _settingsRepository.updateThemeSettings(
        isDarkMode: _isDarkMode,
        useSystemTheme: _useSystemTheme,
      );
    }

    notifyListeners();
  }

  // Light theme
  ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFF1A73E8), // Google Blue
        brightness: Brightness.light,
      ),
      textTheme: GoogleFonts.interTextTheme(ThemeData.light().textTheme),
      appBarTheme: AppBarTheme(
        centerTitle: true,
        elevation: 0,
        backgroundColor: const Color(0xFF1A73E8),
        foregroundColor: Colors.white,
        titleTextStyle: GoogleFonts.inter(
          fontSize: 20,
          fontWeight: FontWeight.w600,
          color: Colors.white,
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        filled: true,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 16,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: Color(0xFF1A73E8),
        foregroundColor: Colors.white,
      ),
      dividerTheme: const DividerThemeData(space: 20, thickness: 1),
    );
  }

  // Dark theme
  ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFF1A73E8), // Google Blue
        brightness: Brightness.dark,
      ),
      textTheme: GoogleFonts.interTextTheme(ThemeData.dark().textTheme),
      appBarTheme: AppBarTheme(
        centerTitle: true,
        elevation: 0,
        backgroundColor: const Color(0xFF1A73E8).withOpacity(0.8),
        foregroundColor: Colors.white,
        titleTextStyle: GoogleFonts.inter(
          fontSize: 20,
          fontWeight: FontWeight.w600,
          color: Colors.white,
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        filled: true,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 16,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: Color(0xFF1A73E8),
        foregroundColor: Colors.white,
      ),
      dividerTheme: const DividerThemeData(space: 20, thickness: 1),
    );
  }
}
