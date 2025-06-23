import json
from rest_framework import generics, permissions
from .serializers import UserSerializer
from rest_framework.response import Response
from rest_framework.views import APIView

from .models import User


class RegisterView(generics.CreateAPIView):
    serializer_class = UserSerializer


class ProtectedView(APIView):
    permission_classes = [permissions.IsAuthenticated]    







