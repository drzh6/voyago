import json
from rest_framework import generics, permissions
from .serializers import UserSerializer
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.tokens import RefreshToken

from .models import User


class RegisterView(generics.CreateAPIView):
    serializer_class = UserSerializer

    def create(self, request, *args, **kwargs):
        response = super().create(request, *args, **kwargs)
        user = User.objects.get(id=response.data['id'])
        refresh = RefreshToken.for_user(user)
        return Response({
            'user': response.data,
            'token': {
                'refresh': str(refresh),
                'access': str(refresh.access_token),
            }
        })


class ProtectedView(APIView):
    permission_classes = [permissions.IsAuthenticated]    







