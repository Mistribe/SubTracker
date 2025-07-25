import 'package:flutter/material.dart';

/// Extension methods for the Color class
extension ColorExtensions on Color {
  /// Creates a copy of this color with the given alpha value.
  /// 
  /// The [alpha] parameter must be between 0 and 1, inclusive.
  Color withValues({double? alpha}) {
    return withOpacity(alpha ?? opacity);
  }
}