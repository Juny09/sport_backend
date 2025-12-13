import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../models/models.dart';
import '../../services/api_service.dart';

class MyBookingsScreen extends StatefulWidget {
  const MyBookingsScreen({super.key});

  @override
  State<MyBookingsScreen> createState() => _MyBookingsScreenState();
}

class _MyBookingsScreenState extends State<MyBookingsScreen> {
  late Future<List<Booking>> _bookingsFuture;

  @override
  void initState() {
    super.initState();
    _refreshBookings();
  }

  void _refreshBookings() {
    setState(() {
      _bookingsFuture = Provider.of<ApiService>(context, listen: false).getMyBookings();
    });
  }

  Future<void> _cancelBooking(int id) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Cancel Booking'),
        content: const Text('Are you sure you want to cancel this booking?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('No')),
          TextButton(onPressed: () => Navigator.pop(context, true), child: const Text('Yes')),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await Provider.of<ApiService>(context, listen: false).cancelBooking(id);
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Booking cancelled.')));
        _refreshBookings();
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Booking>>(
      future: _bookingsFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        } else if (snapshot.hasError) {
          return Center(child: Text('Error: ${snapshot.error}'));
        } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
          return const Center(child: Text('No bookings found.'));
        }

        final bookings = snapshot.data!;
        return ListView.builder(
          itemCount: bookings.length,
          itemBuilder: (context, index) {
            final booking = bookings[index];
            final isCancelled = booking.status == 'cancelled';
            return Card(
              margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: ListTile(
                leading: Icon(
                  isCancelled ? Icons.cancel : Icons.check_circle,
                  color: isCancelled ? Colors.grey : Colors.green,
                ),
                title: Text('Booking #${booking.id}'),
                subtitle: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Unit ID: ${booking.resourceUnitId}'),
                    Text('${_formatDate(booking.startTime)} - ${_formatTime(booking.endTime)}'),
                    Text('Status: ${booking.status}'),
                  ],
                ),
                trailing: (!isCancelled)
                    ? IconButton(
                        icon: const Icon(Icons.delete_outline, color: Colors.red),
                        onPressed: () => _cancelBooking(booking.id),
                      )
                    : null,
              ),
            );
          },
        );
      },
    );
  }

  String _formatDate(String isoTime) {
    final dt = DateTime.parse(isoTime).toLocal();
    return DateFormat('yyyy-MM-dd HH:mm').format(dt);
  }

  String _formatTime(String isoTime) {
    final dt = DateTime.parse(isoTime).toLocal();
    return DateFormat('HH:mm').format(dt);
  }
}
