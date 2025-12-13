import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../models/models.dart';
import '../../services/api_service.dart';

class BookingScreen extends StatefulWidget {
  final Facility facility;

  const BookingScreen({super.key, required this.facility});

  @override
  State<BookingScreen> createState() => _BookingScreenState();
}

class _BookingScreenState extends State<BookingScreen> {
  DateTime _selectedDate = DateTime.now();
  List<UnitAvailability>? _availability;
  bool _isLoading = false;
  
  // Selection
  int? _selectedUnitId;
  String? _selectedStartTime;
  String? _selectedEndTime;
  
  @override
  void initState() {
    super.initState();
    _fetchAvailability();
  }

  Future<void> _fetchAvailability() async {
    setState(() => _isLoading = true);
    try {
      final apiService = Provider.of<ApiService>(context, listen: false);
      final dateStr = DateFormat('yyyy-MM-dd').format(_selectedDate);
      final data = await apiService.getAvailability(
        widget.facility.type,
        dateStr,
        60, // Default duration 60 mins
      );
      setState(() {
        _availability = data;
        _selectedUnitId = null;
        _selectedStartTime = null;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 30)),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
      _fetchAvailability();
    }
  }

  Future<void> _createBooking() async {
    if (_selectedUnitId == null || _selectedStartTime == null || _selectedEndTime == null) return;
    
    setState(() => _isLoading = true);
    try {
      final apiService = Provider.of<ApiService>(context, listen: false);
      await apiService.createBooking(
        _selectedUnitId!,
        _selectedStartTime!,
        _selectedEndTime!,
        'Mobile Booking',
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Booking created successfully!')),
        );
        Navigator.pop(context);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Booking failed: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Book ${widget.facility.name}')),
      body: Column(
        children: [
          ListTile(
            title: Text('Date: ${DateFormat('yyyy-MM-dd').format(_selectedDate)}'),
            trailing: const Icon(Icons.calendar_today),
            onTap: () => _selectDate(context),
          ),
          const Divider(),
          if (_isLoading)
            const Expanded(child: Center(child: CircularProgressIndicator()))
          else if (_availability == null || _availability!.isEmpty)
             const Expanded(child: Center(child: Text('No availability found.')))
          else
            Expanded(
              child: ListView.builder(
                itemCount: _availability!.length,
                itemBuilder: (context, index) {
                  final unit = _availability![index];
                  return Card(
                    margin: const EdgeInsets.all(8),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Text(unit.label, style: const TextStyle(fontWeight: FontWeight.bold)),
                        ),
                        Wrap(
                          spacing: 8,
                          children: unit.freeSlots.map((slot) {
                            final isSelected = _selectedUnitId == unit.unitId && _selectedStartTime == slot.start;
                            return ChoiceChip(
                              label: Text('${_formatTime(slot.start)} - ${_formatTime(slot.end)}'),
                              selected: isSelected,
                              onSelected: (selected) {
                                setState(() {
                                  _selectedUnitId = unit.unitId;
                                  _selectedStartTime = slot.start;
                                  _selectedEndTime = slot.end;
                                });
                              },
                            );
                          }).toList(),
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
        ],
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.all(16.0),
        child: ElevatedButton(
          onPressed: (_selectedUnitId != null && !_isLoading) ? _createBooking : null,
          child: const Text('Confirm Booking'),
        ),
      ),
    );
  }

  String _formatTime(String isoTime) {
    final dt = DateTime.parse(isoTime).toLocal();
    return DateFormat('HH:mm').format(dt);
  }
}
