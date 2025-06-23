from django.db import models


class User(models.Model):
    name = models.CharField(max_length=32)
    surname = models.CharField(max_length=32)
    login = models.CharField(max_length=32, unique=True)
    email = models.CharField(max_length=128, unique=True)
    password = models.CharField(max_length=512)
    avatar = models.CharField(null=True, blank=True)
    is_verified = models.BooleanField(default=False)
    register_date = models.DateTimeField(auto_now_add=True)
    last_seen_date = models.DateTimeField(null=True, blank=True)

    def __str__(self):
        return f"{self.name} {self.surname}"


class Trip(models.Model):
    name = models.CharField(max_length=128)
    description = models.CharField(max_length=2048, null=True, blank=True)
    owner = models.ForeignKey(User, on_delete=models.CASCADE, related_name="trips")
    start_date = models.DateTimeField()
    end_date = models.DateTimeField()
    status = models.CharField(max_length=32)
    is_public = models.BooleanField(default=False)
    invite_code = models.CharField(max_length=8, null=True, blank=True)
    cover_image = models.CharField(max_length=1024, null=True, blank=True, default="")
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.name


class TripParticipant(models.Model):
    trip = models.ForeignKey(Trip, on_delete=models.CASCADE, related_name="participants")
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="trips_participated")
    role = models.CharField(max_length=32)
    joined_at = models.DateTimeField(auto_now_add=True)


class TripPlace(models.Model):
    trip = models.ForeignKey(Trip, on_delete=models.CASCADE, related_name="places")
    name = models.CharField(max_length=64)
    latitude = models.DecimalField(max_digits=9, decimal_places=6)
    longitude = models.DecimalField(max_digits=9, decimal_places=6)
    arrival_date = models.DateTimeField(null=True, blank=True)
    departure_date = models.DateTimeField(null=True, blank=True)


class TripExpense(models.Model):
    trip = models.ForeignKey(Trip, on_delete=models.CASCADE, related_name="expenses")
    name = models.CharField(max_length=32)
    amount = models.IntegerField()
    category = models.CharField(max_length=32)
    paid_by = models.ForeignKey(User, on_delete=models.CASCADE, related_name="expenses_paid")
    created_at = models.DateTimeField(auto_now_add=True)


class TripComment(models.Model):
    trip = models.ForeignKey(Trip, on_delete=models.CASCADE, related_name="comments")
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="comments")
    text = models.TextField()
    created_at = models.DateTimeField(auto_now_add=True)
