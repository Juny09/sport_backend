import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../config.dart';
import '../models/models.dart';

class AuthService extends ChangeNotifier {
  User? _user;
  String? _accessToken;
  bool _isInitializing = true;

  User? get currentUser => _user;
  String? get accessToken => _accessToken;
  bool get isAuthenticated => _accessToken != null;
  bool get isInitializing => _isInitializing;

  Future<void> tryAutoLogin() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      if (!prefs.containsKey('accessToken')) {
        _isInitializing = false;
        notifyListeners();
        return;
      }

      final token = prefs.getString('accessToken');
      final userStr = prefs.getString('userData');

      if (token != null && userStr != null) {
        _accessToken = token;
        _user = User.fromJson(json.decode(userStr));
      }
    } catch (e) {
      // Ignore errors during auto login
    } finally {
      _isInitializing = false;
      notifyListeners();
    }
  }

  Future<void> signIn(String email, String password) async {
    final url = Uri.parse('${Config.apiBaseUrl}/auth/login');
    try {
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'email': email, 'password': password}),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _accessToken = data['access_token'];
        if (data['user'] != null) {
          _user = User.fromJson(data['user']);
        }

        final prefs = await SharedPreferences.getInstance();
        if (_accessToken != null) {
          await prefs.setString('accessToken', _accessToken!);
        }
        if (data['user'] != null) {
          await prefs.setString('userData', json.encode(data['user']));
        }

        notifyListeners();
      } else {
        final error = json.decode(response.body)['error'];
        throw Exception(error ?? 'Login failed');
      }
    } catch (e) {
      rethrow;
    }
  }

  Future<void> signUp(String email, String password) async {
    final url = Uri.parse('${Config.apiBaseUrl}/auth/signup');
    try {
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'email': email, 'password': password}),
      );

      if (response.statusCode == 201) {
        // Signup successful.
        // Optionally login immediately or let user login.
        // For now, we assume user needs to login.
      } else {
        final error = json.decode(response.body)['error'];
        throw Exception(error ?? 'Signup failed');
      }
    } catch (e) {
      rethrow;
    }
  }

  Future<void> signOut() async {
    _accessToken = null;
    _user = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.clear();
    notifyListeners();
  }
}
