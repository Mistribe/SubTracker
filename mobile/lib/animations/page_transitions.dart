import 'package:flutter/material.dart';

/// Custom page transitions following Material Design 3 guidelines
/// 
/// Material Design 3 emphasizes smooth, meaningful transitions between screens.
/// These transitions help users understand the relationship between different parts of the app
/// and provide a more engaging, polished user experience.
class PageTransitions {
  /// Creates a shared axis transition in the X direction (horizontal)
  /// Used for transitions between sibling pages (e.g., tabs)
  static Route<T> horizontalSharedAxisTransition<T>({
    required Widget page,
    RouteSettings? settings,
    bool forward = true,
  }) {
    return PageRouteBuilder<T>(
      settings: settings,
      pageBuilder: (context, animation, secondaryAnimation) => page,
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        // Define the animation curves
        final Curve curve = Curves.easeInOut;
        final Curve reverseCurve = Curves.easeInOut;

        // Create animations for translation and opacity
        final Animatable<Offset> slideTransition = Tween<Offset>(
          begin: Offset(forward ? 1.0 : -1.0, 0.0),
          end: Offset.zero,
        ).chain(CurveTween(curve: curve));

        final Animatable<double> fadeTransition = Tween<double>(
          begin: 0.0,
          end: 1.0,
        ).chain(CurveTween(curve: curve));

        return SlideTransition(
          position: animation.drive(slideTransition),
          child: FadeTransition(
            opacity: animation.drive(fadeTransition),
            child: child,
          ),
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      reverseTransitionDuration: const Duration(milliseconds: 300),
    );
  }

  /// Creates a shared axis transition in the Z direction (depth)
  /// Used for transitions to detail pages or settings
  static Route<T> depthSharedAxisTransition<T>({
    required Widget page,
    RouteSettings? settings,
  }) {
    return PageRouteBuilder<T>(
      settings: settings,
      pageBuilder: (context, animation, secondaryAnimation) => page,
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        // Define the animation curves
        final Curve curve = Curves.easeInOut;

        // Create animations for scale and opacity
        final Animatable<double> scaleTransition = Tween<double>(
          begin: 0.9,
          end: 1.0,
        ).chain(CurveTween(curve: curve));

        final Animatable<double> fadeTransition = Tween<double>(
          begin: 0.0,
          end: 1.0,
        ).chain(CurveTween(curve: curve));

        return ScaleTransition(
          scale: animation.drive(scaleTransition),
          child: FadeTransition(
            opacity: animation.drive(fadeTransition),
            child: child,
          ),
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      reverseTransitionDuration: const Duration(milliseconds: 300),
    );
  }

  /// Creates a vertical shared axis transition (Y-axis)
  /// Used for transitions between hierarchical screens
  static Route<T> verticalSharedAxisTransition<T>({
    required Widget page,
    RouteSettings? settings,
    bool forward = true,
  }) {
    return PageRouteBuilder<T>(
      settings: settings,
      pageBuilder: (context, animation, secondaryAnimation) => page,
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        // Define the animation curves
        final Curve curve = Curves.easeInOut;

        // Create animations for translation and opacity
        final Animatable<Offset> slideTransition = Tween<Offset>(
          begin: Offset(0.0, forward ? 1.0 : -1.0),
          end: Offset.zero,
        ).chain(CurveTween(curve: curve));

        final Animatable<double> fadeTransition = Tween<double>(
          begin: 0.0,
          end: 1.0,
        ).chain(CurveTween(curve: curve));

        return SlideTransition(
          position: animation.drive(slideTransition),
          child: FadeTransition(
            opacity: animation.drive(fadeTransition),
            child: child,
          ),
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      reverseTransitionDuration: const Duration(milliseconds: 300),
    );
  }
}
