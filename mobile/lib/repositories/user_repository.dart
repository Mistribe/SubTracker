import 'package:hive/hive.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';
import '../models/user.dart';
import '../providers/user_provider.dart';

class UserRepository {
  static const String _boxName = 'user_box';
  static const String _currentUserIdKey = 'current_user_id';
  static const String _authTokenKey = 'auth_token';

  late Box<User> _userBox;
  late SharedPreferences _prefs;
  UserProvider? _userProvider;

  // Initialize the repository
  Future<void> initialize() async {
    // Register the adapter if not already registered
    if (!Hive.isAdapterRegistered(5)) {
      Hive.registerAdapter(UserAdapter());
    }

    // Open the box
    _userBox = await Hive.openBox<User>(_boxName);
    _prefs = await SharedPreferences.getInstance();
  }

  // Set the user provider
  void setUserProvider(UserProvider provider) {
    _userProvider = provider;
  }

  // Get the current user
  User? getCurrentUser() {
    final userId = _prefs.getString(_currentUserIdKey);
    if (userId == null) return null;
    return _userBox.get(userId);
  }

  // Get the authentication token
  String? getAuthToken() {
    return _prefs.getString(_authTokenKey);
  }

  // Check if the user is authenticated
  bool isAuthenticated() {
    return getCurrentUser() != null && getAuthToken() != null;
  }

  // Sign up a new user (mock implementation)
  Future<User> signUp(String email, String password, {String? displayName}) async {
    // In a real implementation, this would call an API
    // For now, we'll just create a local user
    final user = User(
      id: const Uuid().v4(),
      email: email,
      displayName: displayName,
    );

    // Save the user to the box
    await _userBox.put(user.id, user);

    // Set as current user
    await _prefs.setString(_currentUserIdKey, user.id);

    // Generate a mock token
    final token = 'mock_token_${const Uuid().v4()}';
    await _prefs.setString(_authTokenKey, token);

    // Notify the provider
    _userProvider?.setCurrentUser(user);

    return user;
  }

  // Sign in a user (mock implementation)
  Future<User?> signIn(String email, String password) async {
    // In a real implementation, this would call an API
    // For now, we'll check if a user with this email exists
    final users = _userBox.values.where((user) => user.email == email);
    if (users.isEmpty) return null;

    final user = users.first;

    // Set as current user
    await _prefs.setString(_currentUserIdKey, user.id);

    // Generate a mock token
    final token = 'mock_token_${const Uuid().v4()}';
    await _prefs.setString(_authTokenKey, token);

    // Notify the provider
    _userProvider?.setCurrentUser(user);

    return user;
  }

  // Sign out the current user
  Future<void> signOut() async {
    await _prefs.remove(_currentUserIdKey);
    await _prefs.remove(_authTokenKey);

    // Notify the provider
    _userProvider?.setCurrentUser(null);
  }

  // Update user information
  Future<User> updateUser(User user) async {
    final updatedUser = user.copyWith(
      updatedAt: DateTime.now(),
    );

    await _userBox.put(user.id, updatedUser);

    // If this is the current user, update the provider
    final currentUserId = _prefs.getString(_currentUserIdKey);
    if (currentUserId == user.id) {
      _userProvider?.setCurrentUser(updatedUser);
    }

    return updatedUser;
  }

  // Delete a user
  Future<void> deleteUser(String userId) async {
    await _userBox.delete(userId);

    // If this is the current user, sign out
    final currentUserId = _prefs.getString(_currentUserIdKey);
    if (currentUserId == userId) {
      await signOut();
    }
  }

  // Close the repository
  Future<void> close() async {
    await _userBox.close();
  }
}