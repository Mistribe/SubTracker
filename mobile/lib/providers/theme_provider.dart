import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';

class ThemeProvider with ChangeNotifier {
  bool _isDarkMode = false;
  bool _useSystemTheme = true;

  ThemeProvider() {
    // Initialize with system theme
    _initializeWithSystemTheme();
  }

  void _initializeWithSystemTheme() {
    var brightness = SchedulerBinding.instance.platformDispatcher.platformBrightness;
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

  void toggleTheme() {
    _useSystemTheme = false;
    _isDarkMode = !_isDarkMode;
    notifyListeners();
  }

  void setSystemTheme() {
    _useSystemTheme = true;
    _initializeWithSystemTheme();
    notifyListeners();
  }

  // Light theme
  ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue,
        brightness: Brightness.light,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      // cardTheme property removed due to type incompatibility
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        filled: true,
      ),
    );
  }

  // Dark theme
  ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue,
        brightness: Brightness.dark,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      // cardTheme property removed due to type incompatibility
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        filled: true,
      ),
    );
  }
}
