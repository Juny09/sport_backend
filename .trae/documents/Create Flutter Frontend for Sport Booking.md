# Goal
Create a **Flutter** frontend in `frontend_sport` that integrates with the existing Go backend. The app will allow users to view facilities, check availability, make bookings, and manage their reservations.

# Tech Stack
- **Framework**: Flutter (Dart)
- **State Management**: `provider` (Simple & effective for this scale)
- **Networking**: `http` (for REST API calls)
- **Authentication**: `supabase_flutter` (to handle Login/Sign up and obtain JWT tokens for the backend)
- **Utilities**: `intl` (Date formatting), `shared_preferences` (Settings/Persistance)

# Project Structure (`frontend_sport/lib`)
- `main.dart`: Entry point, providers setup, route configuration.
- `config.dart`: Environment variables (API URL, Supabase Keys).
- `models/`: Dart data classes (Facility, Unit, Booking, Availability, User).
- `services/`:
  - `auth_service.dart`: Wraps Supabase Auth.
  - `api_service.dart`: HTTP client handling Bearer tokens and backend endpoints.
- `providers/`:
  - `AuthProvider`: Manages user session state.
  - `BookingProvider`: Manages facility list, availability state, and booking actions.
- `screens/`:
  - `auth/login_screen.dart`: Email/Password login.
  - `home/home_screen.dart`: Tab view (Facilities | My Bookings | Profile).
  - `facilities/facilities_list.dart`: List all facilities.
  - `facilities/facility_detail.dart`: View units & check availability.
  - `booking/booking_screen.dart`: Select time slots -> Create Booking.
  - `my_bookings/my_bookings_screen.dart`: List active/past bookings, cancel/reschedule.
  - `admin/admin_dashboard.dart`: (Optional) Quick links for admin actions if role=admin.

# Key Features Implementation
1.  **Authentication**:
    -   Use Supabase to log in.
    -   Pass the JWT (`access_token`) in the `Authorization` header to the Go backend.
2.  **Facilities & Units**:
    -   `GET /facilities` to display cards.
    -   `GET /facilities/:id/units` to show courts/areas.
3.  **Booking Flow**:
    -   User selects a Facility + Date.
    -   App calls `GET /availability` to render a timeline or slot grid.
    -   User selects a free slot -> `POST /bookings`.
4.  **My Bookings**:
    -   `GET /bookings?mine=true`.
    -   Show status tags (Confirmed, Cancelled).
    -   "Cancel" button calls `PATCH .../cancel`.

# Steps
1.  **Initialize**: Run `flutter create frontend_sport`.
2.  **Dependencies**: Add packages to `pubspec.yaml`.
3.  **Scaffold**: Create directory structure and base files.
4.  **Models**: Generate Dart models from Postman JSON.
5.  **Services**: Implement API calls.
6.  **UI**: Build screens with Material Design 3.
7.  **Integration**: Connect UI to Providers/Services.

# Configuration Note
I will create a `lib/config.dart` file where you will need to fill in your **Supabase URL** and **Anon Key**, as well as the **Backend URL** (defaulting to localhost).

Please confirm to proceed with generating the Flutter project.