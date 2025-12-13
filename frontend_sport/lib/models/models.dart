class Facility {
  final int id;
  final String name;
  final String type;
  final bool isActive;

  Facility({
    required this.id,
    required this.name,
    required this.type,
    required this.isActive,
  });

  factory Facility.fromJson(Map<String, dynamic> json) {
    return Facility(
      id: json['ID'],
      name: json['Name'],
      type: json['Type'],
      isActive: json['IsActive'],
    );
  }
}

class User {
  final String id;
  final String email;

  User({required this.id, required this.email});

  factory User.fromJson(Map<String, dynamic> json) {
    return User(id: json['id'], email: json['email']);
  }
}

class Unit {
  final int id;
  final int facilityId;
  final String label;
  final bool isActive;

  Unit({
    required this.id,
    required this.facilityId,
    required this.label,
    required this.isActive,
  });

  factory Unit.fromJson(Map<String, dynamic> json) {
    return Unit(
      id: json['ID'],
      facilityId: json['FacilityID'],
      label: json['Label'],
      isActive: json['IsActive'],
    );
  }
}

class AvailabilitySlot {
  final String start;
  final String end;

  AvailabilitySlot({required this.start, required this.end});

  factory AvailabilitySlot.fromJson(Map<String, dynamic> json) {
    return AvailabilitySlot(start: json['Start'], end: json['End']);
  }
}

class UnitAvailability {
  final int unitId;
  final String label;
  final List<AvailabilitySlot> freeSlots;

  UnitAvailability({
    required this.unitId,
    required this.label,
    required this.freeSlots,
  });

  factory UnitAvailability.fromJson(Map<String, dynamic> json) {
    var slotsList = json['free'] as List;
    List<AvailabilitySlot> slots = slotsList
        .map((i) => AvailabilitySlot.fromJson(i))
        .toList();

    return UnitAvailability(
      unitId: json['unit_id'],
      label: json['label'],
      freeSlots: slots,
    );
  }
}

class Booking {
  final int id;
  final int resourceUnitId;
  final String userId;
  final String startTime;
  final String endTime;
  final String status;
  final double price;

  Booking({
    required this.id,
    required this.resourceUnitId,
    required this.userId,
    required this.startTime,
    required this.endTime,
    required this.status,
    required this.price,
  });

  factory Booking.fromJson(Map<String, dynamic> json) {
    return Booking(
      id: json['ID'],
      resourceUnitId: json['ResourceUnitID'],
      userId: json['UserID'],
      startTime: json['StartTime'],
      endTime: json['EndTime'],
      status: json['Status'],
      price: (json['Price'] as num).toDouble(),
    );
  }
}
