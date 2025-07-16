import 'package:flutter/material.dart';
import '../models/user.dart';
import '../repositories/user_repository.dart';

class UserProvider extends ChangeNotifier {
  final UserRepository _userRepository;
  User? _currentUser;
  bool _isInitialized = false;

  UserProvider({required UserRepository userRepository})
      : _userRepository = userRepository {
    _initialize();
  }

  // Initialize the provider
  Future<void> _initialize() async {
    if (_isInitialized) return;

    // Get the current user from the repository
    _currentUser = _userRepository.getCurrentUser();
    _isInitialized = true;
    notifyListeners();
  }

  // Get the current user
  User? get currentUser => _currentUser;

  // Check if a user is authenticated
  bool get isAuthenticated => _currentUser != null && _userRepository.getAuthToken() != null;

  // Check if the provider is initialized
  bool get isInitialized => _isInitialized;

  // Set the current user (called by the repository)
  void setCurrentUser(User? user) {
    _currentUser = user;
    notifyListeners();
  }

  // Sign up a new user
  Future<User> signUp(String email, String password, {String? displayName}) async {
    final user = await _userRepository.signUp(email, password, displayName: displayName);
    _currentUser = user;
    notifyListeners();
    return user;
  }

  // Sign in a user
  Future<User?> signIn(String email, String password) async {
    final user = await _userRepository.signIn(email, password);
    _currentUser = user;
    notifyListeners();
    return user;
  }

  // Sign out the current user
  Future<void> signOut() async {
    await _userRepository.signOut();
    _currentUser = null;
    notifyListeners();
  }

  // Update user information
  Future<User> updateUser(User user) async {
    final updatedUser = await _userRepository.updateUser(user);
    if (_currentUser?.id == user.id) {
      _currentUser = updatedUser;
      notifyListeners();
    }
    return updatedUser;
  }
}