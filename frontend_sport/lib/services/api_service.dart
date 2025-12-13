import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config.dart';
import '../models/models.dart';

class ApiService {
  String? _authToken;

  void setAuthToken(String? token) {
    _authToken = token;
  }

  Map<String, String> get _headers {
    final headers = {'Content-Type': 'application/json'};
    if (_authToken != null) {
      headers['Authorization'] = 'Bearer $_authToken';
    }
    return headers;
  }

  // Health
  Future<bool> checkHealth() async {
    final response = await http.get(Uri.parse('${Config.apiBaseUrl}/health'));
    return response.statusCode == 200;
  }

  // Facilities
  Future<List<Facility>> getFacilities() async {
    final response = await http.get(
      Uri.parse('${Config.apiBaseUrl}/facilities'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Facility.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load facilities');
    }
  }

  // Units
  Future<List<Unit>> getUnits(int facilityId) async {
    final response = await http.get(
      Uri.parse('${Config.apiBaseUrl}/facilities/$facilityId/units'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Unit.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load units');
    }
  }

  // Availability
  Future<List<UnitAvailability>> getAvailability(
      String facilityType, String date, int duration) async {
    final uri = Uri.parse('${Config.apiBaseUrl}/availability').replace(
        queryParameters: {
          'facility_type': facilityType,
          'date': date,
          'duration': duration.toString()
        });

    final response = await http.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => UnitAvailability.fromJson(json)).toList();
    } else {
      throw Exception('Failed to check availability');
    }
  }

  // Bookings
  Future<void> createBooking(int resourceUnitId, String startTime,
      String endTime, String notes) async {
    final response = await http.post(
      Uri.parse('${Config.apiBaseUrl}/bookings'),
      headers: _headers,
      body: json.encode({
        'resource_unit_id': resourceUnitId,
        'start_time': startTime,
        'end_time': endTime,
        'notes': notes,
      }),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to create booking: ${response.body}');
    }
  }

  Future<List<Booking>> getMyBookings() async {
    final response = await http.get(
      Uri.parse('${Config.apiBaseUrl}/bookings?mine=true'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Booking.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load bookings');
    }
  }

  Future<void> cancelBooking(int bookingId) async {
    final response = await http.patch(
      Uri.parse('${Config.apiBaseUrl}/bookings/$bookingId/cancel'),
      headers: _headers,
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to cancel booking');
    }
  }
}
