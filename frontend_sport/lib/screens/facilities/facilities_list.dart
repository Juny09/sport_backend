import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../services/api_service.dart';
import '../../models/models.dart';
import '../booking/booking_screen.dart';

class FacilitiesListScreen extends StatefulWidget {
  const FacilitiesListScreen({super.key});

  @override
  State<FacilitiesListScreen> createState() => _FacilitiesListScreenState();
}

class _FacilitiesListScreenState extends State<FacilitiesListScreen> {
  late Future<List<Facility>> _facilitiesFuture;

  @override
  void initState() {
    super.initState();
    _facilitiesFuture = Provider.of<ApiService>(context, listen: false).getFacilities();
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Facility>>(
      future: _facilitiesFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        } else if (snapshot.hasError) {
          return Center(child: Text('Error: ${snapshot.error}'));
        } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
          return const Center(child: Text('No facilities found.'));
        }

        final facilities = snapshot.data!;
        return ListView.builder(
          itemCount: facilities.length,
          itemBuilder: (context, index) {
            final facility = facilities[index];
            return Card(
              margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: ListTile(
                leading: Icon(_getIconForType(facility.type)),
                title: Text(facility.name),
                subtitle: Text(facility.type.toUpperCase()),
                trailing: const Icon(Icons.arrow_forward_ios),
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => BookingScreen(facility: facility),
                    ),
                  );
                },
              ),
            );
          },
        );
      },
    );
  }

  IconData _getIconForType(String type) {
    switch (type.toLowerCase()) {
      case 'badminton':
        return Icons.sports_tennis; // Close enough
      case 'tennis':
        return Icons.sports_tennis;
      case 'gym':
        return Icons.fitness_center;
      case 'multipurpose':
        return Icons.business;
      default:
        return Icons.sports;
    }
  }
}
